package proto

import (
	"fmt"

	"github.com/degoke/health-ai-stack/pkg/types"
	"google.golang.org/protobuf/proto"
)

var (
	errNilProto               = fmt.Errorf("proto value is nil")
	errUnsupportedProto       = fmt.Errorf("unsupported Google FHIR R4 proto resource")
	errEmptyContainedResource = fmt.Errorf("contained resource has no resource set")
)

// IsProtoResource reports whether v is a supported Google FHIR R4 protobuf resource.
func IsProtoResource(v any) bool {
	if v == nil {
		return false
	}
	msg, ok := v.(proto.Message)
	if !ok {
		return false
	}
	return isGoogleR4ProtoMessage(msg.ProtoReflect())
}

// ResourceTypeOfProto returns the FHIR resource type for a supported Google FHIR R4 proto value.
func ResourceTypeOfProto(v any) (string, error) {
	if v == nil {
		return "", errNilProto
	}
	msg, ok := v.(proto.Message)
	if !ok {
		return "", errUnsupportedProto
	}
	return resourceTypeFromR4Message(msg.ProtoReflect())
}

func envelopeFromJSON(jsonCodec *types.JSONCodec, resourceType string, jsonBytes []byte, protoVal any) (*types.ResourceEnvelope, error) {
	envelope, err := jsonCodec.ParseJSON(resourceType, jsonBytes)
	if err != nil {
		return nil, err
	}
	envelope.Proto = protoVal
	return envelope, nil
}

func asProtoMessage(v any) (proto.Message, error) {
	if v == nil {
		return nil, errNilProto
	}
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errUnsupportedProto
	}
	if !isGoogleR4ProtoMessage(msg.ProtoReflect()) {
		return nil, errUnsupportedProto
	}
	return msg, nil
}

func assertResourceTypeMatch(expected, actual string) error {
	if expected != "" && actual != expected {
		return fmt.Errorf("resourceType mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}
