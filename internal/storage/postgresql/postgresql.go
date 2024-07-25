package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("%s: unable to parse config: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to connect: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.pool.Ping(ctx)
}

func (s *Storage) Close() {
	s.pool.Close()
}
