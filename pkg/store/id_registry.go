package store

import (
	"context"
	"time"
)

// IDRegistryEntry records an authoritative resource ID registration for one tenant.
type IDRegistryEntry struct {
	ResourceType string    `json:"resourceType"`
	ID           string    `json:"id"`
	RegisteredAt time.Time `json:"registeredAt"`
}

// IDRegistryResult summarizes the outcome of ID registration during a write.
type IDRegistryResult struct {
	ResourceType string `json:"resourceType"`
	ID           string `json:"id"`
	Registered   bool   `json:"registered"`
}

// IDRegistryStore reserves, checks, and registers resource IDs per tenant and resource type.
type IDRegistryStore interface {
	Check(ctx context.Context, resourceType, id string) (bool, error)
	Reserve(ctx context.Context, resourceType, id string) error
	Register(ctx context.Context, resourceType, id string) error
}
