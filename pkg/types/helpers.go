package types

import (
	"encoding/json"
	"fmt"
)

// GetResourceType returns the top-level resourceType from FHIR JSON.
// Returns an error when data is not a JSON object or resourceType is missing or invalid.
func GetResourceType(data []byte) (string, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return "", err
	}
	rt, ok := obj["resourceType"].(string)
	if !ok || rt == "" {
		return "", fmt.Errorf("missing or invalid resourceType")
	}
	return rt, nil
}

// GetID returns the top-level id from FHIR JSON, or "" when id is absent.
func GetID(data []byte) (string, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return "", err
	}
	if v, ok := obj["id"].(string); ok {
		return v, nil
	}
	return "", nil
}

// SetID updates the top-level id (or removes it when id is empty) and returns normalized JSON.
// Unrelated fields are preserved.
func SetID(data []byte, id string) ([]byte, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	if id != "" {
		obj["id"] = id
	} else {
		delete(obj, "id")
	}
	return normalizeObject(obj)
}

// GetMeta extracts meta.versionId and meta.lastUpdated from FHIR JSON.
// Returns zero Meta when meta is absent; errors when meta is not a JSON object.
func GetMeta(data []byte) (*Meta, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	return metaFromObject(obj)
}

// SetMeta updates meta.versionId and meta.lastUpdated and returns normalized JSON.
// Empty VersionID removes versionId; zero LastUpdated removes lastUpdated.
// Removes meta entirely when no fields remain. Unrelated fields are preserved.
func SetMeta(data []byte, meta Meta) ([]byte, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	if err := setMetaOnObject(obj, meta); err != nil {
		return nil, err
	}
	return normalizeObject(obj)
}

// GetReferences recursively collects every reference field in the JSON payload.
// See package documentation for Reference parsing rules.
func GetReferences(data []byte) ([]Reference, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	var refs []Reference
	collectReferences(v, &refs)
	return refs, nil
}

// NormalizeJSON canonicalizes FHIR JSON: valid object, sorted keys, no insignificant whitespace.
// Semantically equivalent payloads with different formatting normalize to identical bytes.
func NormalizeJSON(data []byte) ([]byte, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	return normalizeObject(obj)
}

// HashResource returns the SHA-256 hex digest of the normalized canonical JSON for data.
func HashResource(data []byte) (string, error) {
	normalized, err := NormalizeJSON(data)
	if err != nil {
		return "", err
	}
	return hashBytes(normalized), nil
}
