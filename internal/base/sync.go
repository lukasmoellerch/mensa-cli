package base

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gofrs/flock"
	"github.com/lukasmoellerch/mensa-cli/internal/protobuf/storage"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	base        string
	storagePath string
	fl          *flock.Flock
	Data        storage.Root
	file        *os.File
}

// Create new sync
func NewStore(base string) (*Store, error) {
	lockPath := path.Join(base, "store.lock")
	storagePath := path.Join(base, "store.db")
	fl := flock.New(lockPath)

	return &Store{
		base:        base,
		storagePath: storagePath,
		fl:          fl,
	}, nil
}

func (s *Store) Lock() error {
	return s.fl.Lock()
}

func (s *Store) Unlock() error {
	return s.fl.Unlock()
}

func (s *Store) Open() error {
	f, err := os.OpenFile(s.storagePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Store) Close() error {
	return s.file.Close()
}

func (s *Store) Read() error {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, s.file)
	if err != nil {
		return err
	}
	s.Data.UnmarshalVT(buf.Bytes())
	return nil
}

func (s *Store) Write() error {
	data, err := s.Data.MarshalVT()
	if err != nil {
		return err
	}
	_, err = s.file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) IsEmpty(lang string) bool {
	_, ok := s.Data.LastUpdate[lang]
	return !ok
}

// Initialize sync
func (s *Store) Sync(ctx context.Context, providers []Provider, lang string) error {
	canteens := s.Data.Canteens
	canteensById := make(map[string]*storage.CanteenData, len(canteens))
	seen := make(map[string]bool, len(canteens))
	for _, c := range canteens {
		canteensById[c.Id] = c
	}
	mu := &sync.Mutex{}

	eg, ctx := errgroup.WithContext(ctx)
	for _, provider := range providers {
		provider := provider
		eg.Go(func() error {
			id := provider.Id()
			list, err := provider.FetchCanteens(ctx, lang)
			if err != nil {
				return err
			}
			mu.Lock()
			for _, canteen := range list {
				seen[canteen.ID] = true
				// Augment existing entry
				if c, ok := canteensById[canteen.ID]; ok {
					c.Label[lang] = canteen.Label
				} else {
					c := &storage.CanteenData{
						Id:       canteen.ID,
						Provider: id,
						Meta:     canteen.Meta,
						Label: map[string]string{
							lang: canteen.Label,
						},
					}
					canteensById[c.Id] = c
				}
			}
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	// Remove canteens not seen
	for id := range canteensById {
		if !seen[id] {
			delete(canteensById, id)
		}
	}
	// Build new array of canteens
	canteens = make([]*storage.CanteenData, 0, len(canteensById))
	for _, c := range canteensById {
		canteens = append(canteens, c)
	}
	s.Data.Reset()
	s.Data.Canteens = canteens
	if s.Data.LastUpdate == nil {
		s.Data.LastUpdate = make(map[string]*timestamppb.Timestamp, 1)
	}
	s.Data.LastUpdate[lang] = timestamppb.Now()
	return nil
}

func (s *Store) Filter(ctx context.Context, lang string, keyword string) ([]*storage.CanteenData, error) {
	keyword = strings.ToLower(keyword)

	canteens := []*storage.CanteenData{}
	for _, c := range s.Data.Canteens {
		label, ok := c.Label[lang]
		if !ok {
			continue
		}
		if !strings.Contains(strings.ToLower(label), keyword) {
			continue
		}
		canteens = append(canteens, c)
	}

	return canteens, nil
}
