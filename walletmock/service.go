package main

import "errors"

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Withdraw(playerID string, amount int64) (newBalance int64, err error) {
	p := s.store.GetOrCreatePlayer(playerID)

	if p.Balance < amount {
		return p.Balance, errors.New("insufficient_funds")
	}

	p.Balance -= amount
	return p.Balance, nil
}

func (s *Service) Deposit(playerID string, amount int64) int64 {
	p := s.store.GetOrCreatePlayer(playerID)
	p.Balance += amount
	return p.Balance
}
