package main

type Player struct {
	Balance int64
}

type Store struct {
	Players      map[string]*Player
	Transactions *TransactionLog
}

func NewStore() *Store {
	return &Store{
		Players:      make(map[string]*Player),
		Transactions: NewTransactionLog(),
	}
}

func (s *Store) GetOrCreatePlayer(id string) *Player {
	p, ok := s.Players[id]
	if !ok {
		p = &Player{Balance: 100000}
		s.Players[id] = p
	}
	return p
}
