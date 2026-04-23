package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

type seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) seeder {
	return seeder{
		db: db,
	}
}

func (s seeder) InsertAll() error {
	seed, err := os.ReadFile("postgres/seed/insert_all.sql")

	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(context.Background(), string(seed)); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (s seeder) DeleteAll() error {
	return nil
}
