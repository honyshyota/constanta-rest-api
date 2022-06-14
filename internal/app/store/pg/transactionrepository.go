package pgstore

import (
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
)

type TransactionRepository struct {
	store *Store
}

func (r *TransactionRepository) Create(t *model.Transaction) error {
	if err := t.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO transactions (id, pay, currency, time_create, time_update, trans_status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING trans_id",
		t.UserID,
		t.Pay,
		t.Currency,
		t.TimeCreate,
		t.TimeUpdate,
		t.Status,
	).Scan(&t.TransID)
}
