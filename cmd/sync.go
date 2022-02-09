package cmd

import (
	"context"

	"github.com/lukasmoellerch/mensa-cli/internal/base"
)

func syncStore(ctx context.Context, store *base.Store) error {
	if err := store.Lock(); err != nil {
		return err
	}
	if err := store.Open(); err != nil {
		return err
	}
	if err := store.Read(); err != nil {
		return err
	}
	if store.IsEmpty(langFlag) {
		if err := store.Sync(ctx, providers, langFlag); err != nil {
			return err
		}
		if err := store.Write(); err != nil {
			return err
		}
	}
	if err := store.Close(); err != nil {
		return err
	}

	if err := store.Unlock(); err != nil {
		return err
	}
	return nil
}
