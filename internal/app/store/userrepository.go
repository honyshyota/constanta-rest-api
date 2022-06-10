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
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT user_id, email, time_create, time_update FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.TimeCreate,
		&u.TimeUpdate,
	); err != nil {
		return nil, err
	}

	return u, nil
}
