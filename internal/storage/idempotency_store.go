package storage

import (
	"context"
	"integration-hub/internal/storage/db"
)

type IdempotencyStore struct {
	queries *db.Queries
}

func NewIdempotencyStore(q *db.Queries) *IdempotencyStore {
	return &IdempotencyStore{queries: q}
}

func (s *IdempotencyStore) Get(key string) ([]byte, bool, error) {
	response, err := s.queries.GetIdempotency(context.Background(), key)
	if err != nil {
		return nil, false, nil
	}
	return response, true, nil
}

func (s *IdempotencyStore) Save(key string, response []byte) error {
	return s.queries.SaveIdempotency(context.Background(), db.SaveIdempotencyParams{
		Key:      key,
		Response: response,
	})
}
