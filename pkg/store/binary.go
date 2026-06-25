package store

import (
	"context"
	"time"
)

// BinaryObject is a stored binary payload with lightweight metadata.
type BinaryObject struct {
	Key         string    `json:"key"`
	ContentType string    `json:"contentType,omitempty"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash,omitempty"`
	Data        []byte    `json:"data,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// BinaryStore stores and fetches binary payloads by stable key.
type BinaryStore interface {
	Put(ctx context.Context, obj BinaryObject) error
	Get(ctx context.Context, key string) (*BinaryObject, error)
	Delete(ctx context.Context, key string) error
}
