package store

import (
	"context"
	"time"
)

// AnalyticsEvent is an append-only operational or product analytics record.
type AnalyticsEvent struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Timestamp  time.Time          `json:"timestamp"`
	Dimensions map[string]string  `json:"dimensions,omitempty"`
	Values     map[string]float64 `json:"values,omitempty"`
	Payload    []byte             `json:"payload,omitempty"`
}

// AnalyticsStore appends analytics events and reads them through prepared queries.
// It does not compute aggregates or define reporting semantics.
type AnalyticsStore interface {
	Append(ctx context.Context, event AnalyticsEvent) error
	QueryPrepared(ctx context.Context, query PreparedQuery, args map[string]string) ([]AnalyticsEvent, error)
}
