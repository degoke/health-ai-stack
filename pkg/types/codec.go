package types

import "fmt"

// ResourceCodec parses FHIR JSON into ResourceEnvelope and serializes envelopes back to JSON.
type ResourceCodec interface {
	ParseJSON(resourceType string, data []byte) (*ResourceEnvelope, error)
	ToJSON(resource *ResourceEnvelope) ([]byte, error)
}

// JSONCodec is the default ResourceCodec. It normalizes JSON and derives envelope metadata
// without typed FHIR structs or proto values.
type JSONCodec struct{}

// NewJSONCodec returns a JSONCodec ready for use.
func NewJSONCodec() *JSONCodec {
	return &JSONCodec{}
}

// ParseJSON decodes FHIR JSON, validates resourceType when the argument is non-empty,
// normalizes the payload, and returns a populated ResourceEnvelope with Proto nil.
func (c *JSONCodec) ParseJSON(resourceType string, data []byte) (*ResourceEnvelope, error) {
	obj, err := decodeObject(data)
	if err != nil {
		return nil, err
	}
	envelope, err := envelopeFromObject(obj)
	if err != nil {
		return nil, err
	}
	if resourceType != "" && envelope.ResourceType != resourceType {
		return nil, fmt.Errorf("resourceType mismatch: expected %s, got %s", resourceType, envelope.ResourceType)
	}
	return envelope, nil
}

// ToJSON returns the envelope's normalized JSON bytes. Errors if resource is nil or JSON is empty.
func (c *JSONCodec) ToJSON(resource *ResourceEnvelope) ([]byte, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource is nil")
	}
	if len(resource.JSON) == 0 {
		return nil, fmt.Errorf("resource JSON is empty")
	}
	return resource.JSON, nil
}
