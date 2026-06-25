package proto

// Code generation support for provider registries and adapter boilerplate.
//
// Manual registry construction (see google_r4_registry.go) works for MVP but does not
// scale when adding R5 or alternate providers. Future codegen may emit:
//   - FHIR resource type → oneof field mappings
//   - ProtoCodec wrapper stubs per provider/version
//   - test fixtures for round-trip contract suites shared across R4/R5
//
// Expected inputs: FHIR version, resource type list, provider module path.
// Expected outputs: Go source merged into pkg/proto or a generated/ subdirectory.
//
// No scripts or generated files in MVP — this file documents the intended seam only.

// codegenConfig holds inputs for a future registry generator.
//
//nolint:unused // reserved API surface for future implementation
type codegenConfig struct {
	fhirVersion string
	provider    string
}
