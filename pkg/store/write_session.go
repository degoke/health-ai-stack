package store

import "context"

// WriteSession coordinates transaction-scoped store access for one atomic write.
type WriteSession interface {
	ResourceStore() ResourceStore
	HistoryStore() HistoryStore
	SearchStore() SearchStore
	EventStore() EventStore
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// WriteSessionProvider begins atomic write sessions.
type WriteSessionProvider interface {
	BeginWrite(ctx context.Context) (WriteSession, error)
}
