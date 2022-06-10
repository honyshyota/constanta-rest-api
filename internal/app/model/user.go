package model

import "time"

type User struct {
	TransactionID int
	ID            int
	Email         string
	Pay           float32
	Currency      string
	TimeCreate    time.Time
	TimeUpdate    time.Time
	Status        string
}
