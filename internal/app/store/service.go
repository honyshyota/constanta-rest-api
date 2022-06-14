package store

type Store interface {
	User() UserRepository
	Transaction() TransactionRepository
}
