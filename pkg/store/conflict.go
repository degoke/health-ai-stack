package store

import (
	"context"
	"time"
)

// ConflictRecord captures a write or sync conflict as a first-class persistence record.
type ConflictRecord struct {
	ID              string     `json:"id"`
	ResourceType    string     `json:"resourceType"`
	ResourceID      string     `json:"resourceId"`
	LocalVersionID  string     `json:"localVersionId"`
	RemoteVersionID string     `json:"remoteVersionId"`
	Reason          string     `json:"reason"`
	CreatedAt       time.Time  `json:"createdAt"`
	ResolvedAt      *time.Time `json:"resolvedAt,omitempty"`
}

// ConflictStore records and retrieves conflicts without embedding reconciliation policy.
type ConflictStore interface {
	Append(ctx context.Context, record ConflictRecord) error
	List(ctx context.Context, resourceType, resourceID string) ([]ConflictRecord, error)
	Resolve(ctx context.Context, id string, resolvedAt time.Time) error
}
