package model

import "time"

type User struct {
	ID         int
	Email      string
	Pay        []*Pay
	TimeCreate time.Time
	TimeUpdate time.Time
}

type Pay struct {
	TransactionID int
	Summ          int
	Currency      int
	Status        string
}
