package kv

import (
	"context"

	bolt "github.com/prysmaticlabs/bbolt"
)

type migration func(*bolt.Tx) error

var (
	migrationCompleted = []byte("done")
	upMigrations       []migration
	downMigrations     []migration
)

// RunUpMigrations defined in the upMigrations list.
func (s *Store) RunUpMigrations(ctx context.Context) error {
	// Run any special migrations that require special conditions.
	if err := s.migrateOptimalAttesterProtectionUp(ctx); err != nil {
		return err
	}

	for _, m := range upMigrations {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := s.db.Update(m); err != nil {
			return err
		}
	}
	return nil
}

// RunDownMigrations defined in the downMigrations list.
func (s *Store) RunDownMigrations(ctx context.Context) error {
	// Run any special migrations that require special conditions.
	if err := s.migrateOptimalAttesterProtectionDown(ctx); err != nil {
		return err
	}

	for _, m := range downMigrations {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := s.db.Update(m); err != nil {
			return err
		}
	}
	return nil
}
