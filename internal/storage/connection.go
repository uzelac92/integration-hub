package storage

import (
	"context"
	"integration-hub/config"
	"integration-hub/internal/storage/db"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func Connect(cfg config.Config) *DB {
	pool, err := pgxpool.New(context.Background(), cfg.DbUrl)
	if err != nil {
		log.Fatal("cannot connect to postgres:", err)
	}

	log.Println("connected to postgres")
	return &DB{
		Pool:    pool,
		Queries: db.New(pool),
	}
}
