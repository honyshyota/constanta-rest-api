package store

import "github.com/honyshyota/constanta-rest-api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type TransactionRepository interface {
	Create(*model.Transaction) error
	StatusUpdate(string, int) error
	FindTrans(int) (*model.Transaction, error)
}
