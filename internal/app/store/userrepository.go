package store

import (
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, time_create, time_update) VALUES ($1, $2, $3) RETURNING user_id",
		u.Email,
		u.TimeCreate,
		u.TimeUpdate,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	// u := &model.User{}
	// if err := r.store.db.QueryRow("SELECT transaction_id, user_id, email, pay, currency, time_create, time")

	return nil, nil
}
