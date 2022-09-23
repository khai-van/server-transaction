package transaction

import (
	"app/pkg/database"
	"app/pkg/sync"
)

type Repository interface { // Repo methods here...
	GetUser(account string) (User, error)
	CreateUser(account string) (User, error)

	Transfer(from string, to string, amount uint64) (Transaction, error)
	GetListTransaction(account string) ([]Transaction, error)
}

type repo struct { // Hold database instance
	db *database.DB
	wg *sync.WaitMapObject // lock process by account name
}

func NewRepository(db *database.DB) Repository {
	return &repo{db: db, wg: sync.WaitMap()}
}
