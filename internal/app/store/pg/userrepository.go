package pgstore

import (
	"database/sql"

	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (id, email, pay, currency, time_create, time_update, transaction_status, encrypted_password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING transaction_id",
		u.ID,
		u.Email,
		u.Pay,
		u.Currency,
		u.TimeCreate,
		u.TimeUpdate,
		u.Status,
		u.EncryptedPassword,
	).Scan(&u.TransactionID)

}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT transaction_id, id, email, pay, currency, time_create, time_update, transaction_status, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.TransactionID,
		&u.ID,
		&u.Email,
		&u.Pay,
		&u.Currency,
		&u.TimeCreate,
		&u.TimeUpdate,
		&u.Status,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
