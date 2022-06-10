package model

import "time"

type User struct {
	ID         int       `pg:"user_id"`
	Email      string    `pg:"email"`
	Pay        []*Pay    `pg:"pay"`
	TimeCreate time.Time `pg:"time_create"`
	TimeUpdate time.Time `pg:"time_updata"`
}

type Pay struct {
	TransactionID int    `pg:"transaction_id"`
	Summ          int    `pg:"summ"`
	Currency      int    `pg:"currency"`
	Status        string `pg:"transaction_status"`
}
