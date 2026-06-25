package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func decodeObject(data []byte) (map[string]interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	obj, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected JSON object")
	}
	return obj, nil
}

func marshalCanonical(obj map[string]interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func normalizeObject(obj map[string]interface{}) ([]byte, error) {
	return marshalCanonical(obj)
}

func parseFHIRTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range formats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid FHIR timestamp: %s", s)
}

func formatFHIRTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func parseReference(raw string) Reference {
	ref := Reference{Raw: raw}
	if raw == "" {
		return ref
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") ||
		strings.HasPrefix(raw, "urn:") || strings.HasPrefix(raw, "#") {
		return ref
	}
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" && !strings.Contains(parts[0], ":") {
		ref.ResourceType = parts[0]
		ref.ID = parts[1]
	}
	return ref
}

func collectReferences(v interface{}, refs *[]Reference) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			if k == "reference" {
				if s, ok := child.(string); ok {
					*refs = append(*refs, parseReference(s))
				}
			} else {
				collectReferences(child, refs)
			}
		}
	case []interface{}:
		for _, item := range val {
			collectReferences(item, refs)
		}
	}
}

func metaFromObject(obj map[string]interface{}) (*Meta, error) {
	metaRaw, ok := obj["meta"]
	if !ok {
		return &Meta{}, nil
	}
	metaMap, ok := metaRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("meta must be a JSON object")
	}
	meta := &Meta{}
	if v, ok := metaMap["versionId"]; ok {
		if s, ok := v.(string); ok {
			meta.VersionID = s
		}
	}
	if v, ok := metaMap["lastUpdated"]; ok {
		if s, ok := v.(string); ok {
			t, err := parseFHIRTime(s)
			if err != nil {
				return nil, err
			}
			meta.LastUpdated = t
		}
	}
	return meta, nil
}

func setMetaOnObject(obj map[string]interface{}, meta Meta) error {
	metaMap, ok := obj["meta"].(map[string]interface{})
	if !ok {
		metaMap = make(map[string]interface{})
		obj["meta"] = metaMap
	}
	if meta.VersionID != "" {
		metaMap["versionId"] = meta.VersionID
	} else {
		delete(metaMap, "versionId")
	}
	if !meta.LastUpdated.IsZero() {
		metaMap["lastUpdated"] = formatFHIRTime(meta.LastUpdated)
	} else {
		delete(metaMap, "lastUpdated")
	}
	if len(metaMap) == 0 {
		delete(obj, "meta")
	}
	return nil
}

func envelopeFromObject(obj map[string]interface{}) (*ResourceEnvelope, error) {
	rt, ok := obj["resourceType"].(string)
	if !ok || rt == "" {
		return nil, fmt.Errorf("missing or invalid resourceType")
	}
	normalized, err := normalizeObject(obj)
	if err != nil {
		return nil, err
	}
	meta, err := metaFromObject(obj)
	if err != nil {
		return nil, err
	}
	id := ""
	if v, ok := obj["id"].(string); ok {
		id = v
	}
	return &ResourceEnvelope{
		ResourceType: rt,
		ID:           id,
		VersionID:    meta.VersionID,
		LastUpdated:  meta.LastUpdated,
		JSON:         normalized,
		Hash:         hashBytes(normalized),
	}, nil
}
