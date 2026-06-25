package types

import "time"

// Reference is a parsed FHIR Reference.reference value.
// Raw always holds the original string; ResourceType and ID are set only for typed
// relative references (see package documentation for parsing rules).
type Reference struct {
	ResourceType string
	ID           string
	Raw          string
}

// Meta holds lightweight FHIR meta fields extracted from or applied to resource JSON.
// It is not a full FHIR Meta complex type — only versionId and lastUpdated are supported.
type Meta struct {
	VersionID   string
	LastUpdated time.Time
}

// Identifier is a lightweight FHIR Identifier for use in future helpers and callers.
// MVP exposes the type only; no JSON extraction helpers are provided yet.
type Identifier struct {
	System string
	Value  string
	Use    string
	Type   string
}
