package store

import (
	"context"
	"time"
)

// Cursor is a named checkpoint for event consumers, sync workers, or indexers.
type Cursor struct {
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CursorStore persists consumer progress independently from event storage.
type CursorStore interface {
	GetCursor(ctx context.Context, name string) (*Cursor, error)
	UpsertCursor(ctx context.Context, cursor Cursor) error
	DeleteCursor(ctx context.Context, name string) error
}
