package model

import (
	"testing"
)

func TestUser(t *testing.T) *User {
	return &User{
		ID:       1,
		Password: "password",
		Email:    "user@example.com",
	}
}
