package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	TransactionID     int       `json:"transaction_id"`
	ID                int       `json:"id"`
	Password          string    `json:"password,omitempty"`
	Email             string    `json:"email"`
	Pay               float32   `json:"pay"`
	Currency          string    `json:"currency"`
	TimeCreate        time.Time `json:"time_create"`
	TimeUpdate        time.Time `json:"time_update"`
	Status            string    `json:"status"`
	EncryptedPassword string    `json:"-"`
}

func (u *User) Validate() error {
	if u.Pay <= 0 {
		u.Status = "error"
	}

	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
		validation.Field(&u.Currency, is.CurrencyCode),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptedString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func encryptedString(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
