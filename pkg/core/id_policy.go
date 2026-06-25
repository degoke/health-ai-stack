package core

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

var fhirIDPattern = regexp.MustCompile(`^[A-Za-z0-9\-\.]{1,64}$`)

// ResourceIDPolicy validates and generates FHIR resource IDs.
type ResourceIDPolicy interface {
	// Validate checks that id is acceptable for resourceType.
	Validate(resourceType, id string) error
	// Generate returns a new id when the caller omits one on create.
	Generate(resourceType string) (string, error)
}

// DefaultIDPolicy accepts FHIR-constrained IDs and generates UUIDs when omitted.
type DefaultIDPolicy struct{}

// Validate checks that id satisfies FHIR id syntax.
func (DefaultIDPolicy) Validate(_ string, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if !fhirIDPattern.MatchString(id) {
		return fmt.Errorf("id %q does not match FHIR id syntax", id)
	}
	return nil
}

// Generate returns a new UUID suitable for use as a FHIR resource id.
func (DefaultIDPolicy) Generate(_ string) (string, error) {
	return uuid.NewString(), nil
}
