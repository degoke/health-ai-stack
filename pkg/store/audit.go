package store

import (
	"context"
	"time"
)

// AuditRecord is an append-only compliance or security audit entry.
type AuditRecord struct {
	ID           string            `json:"id"`
	Timestamp    time.Time         `json:"timestamp"`
	Actor        string            `json:"actor"`
	Action       string            `json:"action"`
	ResourceType string            `json:"resourceType,omitempty"`
	ResourceID   string            `json:"resourceId,omitempty"`
	Outcome      string            `json:"outcome,omitempty"`
	Details      map[string]string `json:"details,omitempty"`
}

// AuditQuery selects audit records with simple lookup filters.
type AuditQuery struct {
	ResourceType string    `json:"resourceType,omitempty"`
	ResourceID   string    `json:"resourceId,omitempty"`
	Actor        string    `json:"actor,omitempty"`
	After        time.Time `json:"after,omitempty"`
	Before       time.Time `json:"before,omitempty"`
	Limit        int       `json:"limit,omitempty"`
}

// AuditStore appends audit records and lists them for review pipelines.
type AuditStore interface {
	Append(ctx context.Context, record AuditRecord) error
	List(ctx context.Context, query AuditQuery) ([]AuditRecord, error)
}
