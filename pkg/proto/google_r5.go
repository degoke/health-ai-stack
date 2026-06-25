package proto

// googleR5Codec will implement ProtoCodec for Google FHIR Go R5.
//
// Planned shape (mirror google_r4.go):
//   - type GoogleR5Codec struct
//   - func NewGoogleR5Codec() *GoogleR5Codec
//   - jsonformat.NewUnmarshaller / NewMarshaller with fhirversion.R5
//   - r5 resource registry (buildR5ResourceRegistry from R5 ContainedResource oneof)
//   - ParseJSONToProto returns R5 *ContainedResource as any
//
// IsProtoResource and ResourceTypeOfProto will need provider dispatch once multiple
// codecs are active (e.g. detect R4 vs R5 message descriptors).
//
// Not implemented in MVP.
//
//nolint:unused // reserved API surface for future implementation
type googleR5Codec struct{}

// compile-time check once implemented:
// var _ ProtoCodec = (*GoogleR5Codec)(nil)
