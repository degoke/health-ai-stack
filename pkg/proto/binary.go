package proto

// protoBinaryCodec marshals and unmarshals FHIR protobuf wire format.
//
// Intended for optional companion storage (e.g. resource_proto_blob) alongside
// canonical resource_json — never as a replacement for JSON storage.
//
// Planned methods:
//   - MarshalProto(proto any) ([]byte, error)
//   - UnmarshalProto(resourceType string, data []byte) (any, error)
//
// Implementation will delegate to google.golang.org/protobuf/proto and reuse the
// same provider registry as JSON codecs to select the correct message type.
//
// Not implemented in MVP.
//
//nolint:unused // reserved API surface for future implementation
type protoBinaryCodec interface {
	marshalProto(proto any) ([]byte, error)
	unmarshalProto(resourceType string, data []byte) (any, error)
}
