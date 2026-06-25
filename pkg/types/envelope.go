package types

import "time"

// ResourceEnvelope is the generic runtime container for a FHIR resource.
// JSON and Hash are always canonical; Proto is optional (set by pkg/proto on proto paths).
type ResourceEnvelope struct {
	ResourceType string
	ID           string
	VersionID    string
	LastUpdated  time.Time
	JSON         []byte
	Proto        any
	Hash         string
}
