package models

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	PGDATABASE = "hotel-system"
	PGURI      = "postgres://myuser:mypassword@localhost:5432/hotel-system"
)

var (
	PGInstance *PostgresInstance
	pgOnce     sync.Once
	Ctx        = context.Background()
)

type PostgresInstance struct {
	DB *pgxpool.Pool
}

func (pg *PostgresInstance) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *PostgresInstance) Close() {
	pg.DB.Close()
}

func NewPostgresInstance(ctx context.Context, connString string) *PostgresInstance {
	pgOnce.Do(func() {
		pool, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatal(err)
		}
		PGInstance = &PostgresInstance{
			DB: pool,
		}
	})
	return PGInstance
}
