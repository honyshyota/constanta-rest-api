package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		ID:       1,
		Password: "password",
		Email:    "user@example.com",
	}
}

func TestTransaction(t *testing.T) *Transaction {
	return &Transaction{
		TransID:    1,
		UserID:     1,
		Email:      "user@example.com",
		Pay:        100.00,
		Currency:   "RUB",
		TimeCreate: time.Now(),
		TimeUpdate: time.Now(),
		Status:     "New",
	}
}
