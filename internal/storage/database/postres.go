package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(databaseURL string) error {
	pool, err := newPool(databaseURL)
	if err != nil {
		return err
	}

	if err = migrate(pool, 1); err != nil {
		return fmt.Errorf("postgres migration error: %w", err)
	}

	pool.Close()

	return nil
}

func newPool(databaseURL string) (*pgxpool.Pool, error) {
	newCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	pool, err := pgxpool.New(newCtx, databaseURL)
	if err != nil {
		return pool, fmt.Errorf("postgres connection error: %w", err)
	}

	return pool, nil
}

type StoragePostgres struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func New(databaseURL string) (*StoragePostgres, error) {
	pool, err := newPool(databaseURL)
	if err != nil {
		return nil, err
	}
	return &StoragePostgres{
		db:      pool,
		timeout: time.Second * 2,
	}, nil
}
