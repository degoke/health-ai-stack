package store

import "context"

// Transaction commits or rolls back a unit of work for transaction-capable backends.
type Transaction interface {
	Commit() error
	Rollback() error
}

// Transactor begins transactional write scopes.
// Only backends that support atomic writes need to implement this interface.
//
// When BeginTx is supported, all writes performed through the transaction-scoped
// store or session must commit atomically according to that backend's rules.
type Transactor interface {
	BeginTx(ctx context.Context) (Transaction, error)
}
