package teststore

import (
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
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

func (r *TransactionRepository) FindTrans(id int) (*model.Transaction, error) {
	transaction, ok := r.transactions[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return transaction, nil
}

func (r *TransactionRepository) StatusUpdate(status string, id int) error {
	return nil
}

func (r *TransactionRepository) Find(data string) ([]*model.Transaction, error) {
	return nil, nil
}

func (r *TransactionRepository) Delete(id int) error {
	return nil
}
