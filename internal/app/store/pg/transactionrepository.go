package pgstore

import (
	"database/sql"
	"errors"

	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
)

type TransactionRepository struct {
	store *Store
}

func (r *TransactionRepository) Create(t *model.Transaction) error {
	if err := t.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO transactions (id, email, pay, currency, time_create, time_update, trans_status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING trans_id",
		t.UserID,
		t.Email,
		t.Pay,
		t.Currency,
		t.TimeCreate,
		t.TimeUpdate,
		t.Status,
	).Scan(&t.TransID)
}

func (r *TransactionRepository) StatusUpdate(status string, id int) error {
	if status != "success" && status != "failure" {
		return errors.New("incorrect input")
	}

	if _, err := r.store.db.Exec(
		"UPDATE transactions SET trans_status = $1 WHERE trans_id = $2",
		status,
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (r *TransactionRepository) FindTrans(id int) (*model.Transaction, error) {
	trans := &model.Transaction{}
	if err := r.store.db.QueryRow(
		"SELECT trans_status FROM transactions WHERE trans_id = $1",
		id,
	).Scan(
		&trans.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return trans, nil
}
