package pgstore

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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
		"UPDATE transactions SET trans_status = $1, time_update = $2 WHERE trans_id = $3",
		status,
		time.Now(),
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
		"SELECT id, trans_status FROM transactions WHERE trans_id = $1",
		id,
	).Scan(
		&trans.UserID,
		&trans.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return trans, nil
}

func (r *TransactionRepository) Find(data string) ([]*model.Transaction, error) {
	var transResult []*model.Transaction

	if err := validation.Validate(data, validation.Required, is.Email); err != nil {
		id, err := strconv.Atoi(data)
		if err != nil {
			return nil, err
		}

		rows, err := r.store.db.Query(
			"SELECT trans_id, id, email, pay, currency, time_create, time_update, trans_status FROM transactions WHERE id = $1",
			id,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			result := &model.Transaction{}
			if err := rows.Scan(
				&result.TransID,
				&result.UserID,
				&result.Email,
				&result.Pay,
				&result.Currency,
				&result.TimeCreate,
				&result.TimeUpdate,
				&result.Status,
			); err != nil {
				return nil, err
			}
			transResult = append(transResult, result)
		}
	} else {
		rows, err := r.store.db.Query(
			"SELECT trans_id, id, email, pay, currency, time_create, time_update, trans_status FROM transactions WHERE email = $1",
			data,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}

			return nil, err
		}

		defer rows.Close()

		fmt.Println(rows.Columns())

		for rows.Next() {
			result := &model.Transaction{}
			if err := rows.Scan(
				&result.TransID,
				&result.UserID,
				&result.Email,
				&result.Pay,
				&result.Currency,
				&result.TimeCreate,
				&result.TimeUpdate,
				&result.Status,
			); err != nil {
				return nil, err
			}
			transResult = append(transResult, result)
		}
	}
	return transResult, nil
}

func (r *TransactionRepository) Delete(id int) error {
	if _, err := r.store.db.Exec(
		"DELETE FROM transactions WHERE trans_id = $1",
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	return nil
}
