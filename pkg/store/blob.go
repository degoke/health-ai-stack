package store

import (
	"context"
	"time"
)

// BlobObject is a stored blob payload or opaque backend reference.
// BinaryStore is intended for small inline payloads; BlobStore covers larger
// or externally referenced content without exposing object-storage SDK details.
type BlobObject struct {
	Key         string    `json:"key"`
	ContentType string    `json:"contentType,omitempty"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash,omitempty"`
	Data        []byte    `json:"data,omitempty"`
	Location    string    `json:"location,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// BlobStore stores and fetches blob payloads by stable key.
type BlobStore interface {
	Put(ctx context.Context, obj BlobObject) error
	Get(ctx context.Context, key string) (*BlobObject, error)
	Head(ctx context.Context, key string) (*BlobObject, error)
	Delete(ctx context.Context, key string) error
}
