package store

type Store interface {
	User() UserRepository
}

type TransactionStore interface {
	Transaction() TransactionRepository
}
