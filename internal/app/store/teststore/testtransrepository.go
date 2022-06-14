package teststore

import (
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
)

type TransactionRepository struct {
	store        *Store
	transactions map[int]*model.Transaction
}

func (r *TransactionRepository) Create(t *model.Transaction) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.TransID = len(r.transactions) + 1
	r.transactions[t.TransID] = t

	return nil
}
