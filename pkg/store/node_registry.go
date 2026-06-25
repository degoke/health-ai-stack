package store

import (
	"context"
	"time"
)

// NodeRecord describes a registered edge or cloud sync node.
type NodeRecord struct {
	NodeID       string            `json:"nodeId"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	RegisteredAt time.Time         `json:"registeredAt"`
}

// NodeRegistryStore persists node registration metadata for sync coordination.
type NodeRegistryStore interface {
	Register(ctx context.Context, node NodeRecord) error
	Get(ctx context.Context, nodeID string) (*NodeRecord, error)
	List(ctx context.Context) ([]NodeRecord, error)
	Unregister(ctx context.Context, nodeID string) error
}
