package model

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Transaction struct {
	TransID    int       `json:"trans_id"`
	UserID     int       `json:"user_id"`
	Email      string    `json:"email"`
	Pay        float32   `json:"pay"`
	Currency   string    `json:"currency"`
	TimeCreate time.Time `json:"time_create"`
	TimeUpdate time.Time `json:"time_update"`
	Status     string    `json:"status"`
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
