package proto

// protoValidator performs provider-aware structural validation on parsed protos.
//
// JSON parsing via jsonformat already runs Google's extended validation when using
// NewUnmarshaller. This seam is for additional typed checks beyond JSON ingest:
//   - required fields for domain-specific profiles
//   - cross-field constraints
//   - reference integrity against a local index
//
// Planned shape:
//   - func ValidateProto(resourceType string, proto any) error
//   - provider-specific validators registered alongside each ProtoCodec
//
// Not implemented in MVP.
//
//nolint:unused // reserved API surface for future implementation
type protoValidator interface {
	validateProto(resourceType string, proto any) error
}
