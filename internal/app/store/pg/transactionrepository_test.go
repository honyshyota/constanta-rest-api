package pgstore_test

import (
	"testing"

	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	pgstore "github.com/honyshyota/constanta-rest-api/internal/app/store/pg"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Create(t *testing.T) {
	db, teardown := pgstore.TestDB(t, databaseURL)
	defer teardown("transactions")

	store := pgstore.New(db)
	transaction := model.TestTransaction(t)
	assert.NoError(t, store.Transaction().Create(transaction))
	assert.NotNil(t, transaction)
}
