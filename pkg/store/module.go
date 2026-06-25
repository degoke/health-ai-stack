package store

import (
	"context"
	"time"
)

// ModuleRecord describes a registered runtime module or extension.
type ModuleRecord struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	RegisteredAt time.Time         `json:"registeredAt"`
}

// ModuleStore persists module registration metadata without loading code or
// enforcing activation policy.
type ModuleStore interface {
	Register(ctx context.Context, module ModuleRecord) error
	Get(ctx context.Context, name string) (*ModuleRecord, error)
	List(ctx context.Context) ([]ModuleRecord, error)
	Unregister(ctx context.Context, name string) error
}
