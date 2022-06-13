package pgstore

import (
	"database/sql"

	"github.com/honyshyota/constanta-rest-api/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db                    *sql.DB
	UserRepository        *UserRepository
	TransactionRepository *TransactionRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}

func (s *Store) Transaction() store.TransactionRepository {
	if s.TransactionRepository != nil {
		return s.TransactionRepository
	}

	s.TransactionRepository = &TransactionRepository{
		store: s,
	}

	return s.TransactionRepository
}
