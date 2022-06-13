package pgstore_test

import (
	"testing"

	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
	pgstore "github.com/honyshyota/constanta-rest-api/internal/app/store/pg"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := pgstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := pgstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := pgstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := pgstore.New(db)
	u1 := model.TestUser(t)
	s.User().Create(u1)
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := pgstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := pgstore.New(db)
	email := "user@example.com"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	s.User().Create(u)
	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
