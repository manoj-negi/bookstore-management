package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connString string) (*Postgres, error) {
	db, err := connectToDB(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create db connection pool: %w", err)
	}
	fmt.Println("Database connected")
	return &Postgres{DB: db}, nil
}

func connectToDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		// Check for auth error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "28P01" {
				return nil, fmt.Errorf("invalid username/password")
			}
		}
		return nil, err
	}

	return db, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.DB.Ping(ctx)
}

func (p *Postgres) Close() error {
	defer p.DB.Close()
	return nil
}
