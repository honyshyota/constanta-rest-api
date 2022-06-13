package model

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Transaction struct {
	TransID    int
	UserID     int
	Email      string
	Pay        float32
	Currency   string
	TimeCreate time.Time
	TimeUpdate time.Time
	Status     string
}

func (t *Transaction) Validate() error {
	if t.Pay < 0 {
		return errors.New("payment must be more than 0")
	}

	return validation.ValidateStruct(
		t,
		validation.Field(&t.Pay, validation.Required),
		validation.Field(&t.Currency, validation.Required, is.CurrencyCode),
	)
}
