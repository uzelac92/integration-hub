package storage

import (
	"context"
	"integration-hub/internal/storage/db"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func Connect() *DB {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://hub:hubpass@localhost:5432/hubdb"
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("cannot connect to postgres:", err)
	}

	log.Println("connected to postgres")

	return &DB{
		Pool:    pool,
		Queries: db.New(pool),
	}
}
