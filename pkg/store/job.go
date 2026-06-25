package store

import (
	"context"
	"time"
)

// JobStatus describes persisted background job state.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// JobRecord is a durable background job entry.
type JobRecord struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   []byte    `json:"payload,omitempty"`
	Status    JobStatus `json:"status"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RunAfter  time.Time `json:"runAfter,omitempty"`
	LastError string    `json:"lastError,omitempty"`
}

// JobStore persists background jobs for workers to enqueue, claim, and update.
// Scheduling, retry, and execution policy remain outside this package.
type JobStore interface {
	Enqueue(ctx context.Context, job JobRecord) error
	ClaimNext(ctx context.Context, jobType string) (*JobRecord, error)
	Update(ctx context.Context, job JobRecord) error
	Get(ctx context.Context, id string) (*JobRecord, error)
}
