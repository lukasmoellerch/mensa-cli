package zuerichmensa

import (
	"bytes"
	"context"
	"strings"

	"encoding/gob"

	"git.mills.io/prologic/bitcask"
)

type FacilityStore struct {
	db *bitcask.Bitcask
}

// Create new sync
func NewStore(base string) (*FacilityStore, error) {
	db, err := bitcask.Open(base)
	if err != nil {
		return nil, err
	}
	return &FacilityStore{
		db: db,
	}, nil
}

func getValue(value Facility) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *FacilityStore) IsEmpty() bool {
	return s.db.Len() == 0
}

// Initialize sync
func (s *FacilityStore) Sync(ctx context.Context) error {
	err := s.db.DeleteAll()
	if err != nil {
		return err
	}
	facilities, err := FetchFacilitiesEth(ctx)
	if err != nil {
		return err
	}
	for _, facility := range facilities.Facilites {
		value, err := getValue(facility)
		if err != nil {
			return err
		}
		key := []byte(strings.ToLower(facility.Label))
		err = s.db.Put(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FacilityStore) Filter(ctx context.Context, keyword string) ([]Facility, error) {
	facilities := make([]Facility, 0)
	keyword = strings.ToLower(keyword)
	err := s.db.Fold(func(key []byte) error {
		keyString := string(key)
		if !strings.Contains(keyString, keyword) {
			return nil
		}

		value, err := s.db.Get(key)
		if err != nil {
			return err
		}
		var facility Facility
		err = gob.NewDecoder(bytes.NewReader(value)).Decode(&facility)
		if err != nil {
			return err
		}
		facilities = append(facilities, facility)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return facilities, nil
}
