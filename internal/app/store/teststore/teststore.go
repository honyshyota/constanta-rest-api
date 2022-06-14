package teststore

import (
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
)

type Store struct {
	UserRepository        *UserRepository
	TransactionRepository *TransactionRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.UserRepository
}

func (s *Store) Transaction() store.TransactionRepository {
	if s.TransactionRepository != nil {
		return s.TransactionRepository
	}

	s.TransactionRepository = &TransactionRepository{
		store:        s,
		transactions: make(map[int]*model.Transaction),
	}

	return s.TransactionRepository
}
