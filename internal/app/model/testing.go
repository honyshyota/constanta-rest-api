package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		ID: 1,
		Password: "password",
		Email: "user@example.com",
		Pay: 1.0,
		Currency: "RUB",
		TimeCreate: time.Now(),
		TimeUpdate: time.Now(),
		Status: "Succes",
		
	}
}
