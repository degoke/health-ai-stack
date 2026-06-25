package store

import (
	"context"
	"time"
)

// MaterializedViewRecord is a persisted projection entry for a named view.
type MaterializedViewRecord struct {
	ViewName  string    `json:"viewName"`
	Key       string    `json:"key"`
	Payload   []byte    `json:"payload,omitempty"`
	Version   int64     `json:"version,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MaterializedViewStore persists read-optimized projections independently from
// current-state resource storage. Callers own projection logic and refresh policy.
type MaterializedViewStore interface {
	Upsert(ctx context.Context, record MaterializedViewRecord) error
	Get(ctx context.Context, viewName, key string) (*MaterializedViewRecord, error)
	Delete(ctx context.Context, viewName, key string) error
	ListKeys(ctx context.Context, viewName string) ([]string, error)
}
