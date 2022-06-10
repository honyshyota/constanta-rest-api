package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	TransactionID     int
	ID                int
	Password          string
	Email             string
	Pay               float32
	Currency          string
	TimeCreate        time.Time
	TimeUpdate        time.Time
	Status            string
	EncryptedPassword string
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

func encryptedString(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
