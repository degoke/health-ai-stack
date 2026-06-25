package proto

import (
	"unicode"

	rpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const r4ContainedResourceOneof = "oneof_resource"

type r4ResourceEntry struct {
	fhirType    string
	oneofField  protoreflect.Name
	messageName protoreflect.FullName
}

type r4ResourceRegistry struct {
	byFHIRType map[string]r4ResourceEntry
	byMessage  map[protoreflect.FullName]string
}

var defaultR4Registry = buildR4ResourceRegistry()

func buildR4ResourceRegistry() *r4ResourceRegistry {
	reg := &r4ResourceRegistry{
		byFHIRType: make(map[string]r4ResourceEntry),
		byMessage:  make(map[protoreflect.FullName]string),
	}

	cr := (&rpb.ContainedResource{}).ProtoReflect()
	od := cr.Descriptor().Oneofs().ByName(r4ContainedResourceOneof)
	if od == nil {
		return reg
	}

	for i := 0; i < od.Fields().Len(); i++ {
		fd := od.Fields().Get(i)
		fhirType := fhirTypeFromOneofField(fd)
		entry := r4ResourceEntry{
			fhirType:    fhirType,
			oneofField:  fd.Name(),
			messageName: fd.Message().FullName(),
		}
		reg.byFHIRType[fhirType] = entry
		reg.byMessage[fd.Message().FullName()] = fhirType
	}

	return reg
}

func fhirTypeFromOneofField(fd protoreflect.FieldDescriptor) string {
	if jsonName := fd.JSONName(); jsonName != "" {
		return camelToPascal(jsonName)
	}
	return camelToPascal(string(fd.Name()))
}

func camelToPascal(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (reg *r4ResourceRegistry) fhirTypeForMessage(name protoreflect.FullName) (string, bool) {
	rt, ok := reg.byMessage[name]
	return rt, ok
}

func (reg *r4ResourceRegistry) isKnownMessage(name protoreflect.FullName) bool {
	_, ok := reg.byMessage[name]
	return ok
}

func (reg *r4ResourceRegistry) isContainedResource(name protoreflect.FullName) bool {
	return name == containedResourceFullName()
}

func containedResourceFullName() protoreflect.FullName {
	return (&rpb.ContainedResource{}).ProtoReflect().Descriptor().FullName()
}

func resourceTypeFromR4ContainedResource(cr *rpb.ContainedResource) (string, error) {
	if cr == nil {
		return "", errNilProto
	}
	msg := cr.ProtoReflect()
	od := msg.Descriptor().Oneofs().ByName(r4ContainedResourceOneof)
	if od == nil {
		return "", errUnsupportedProto
	}
	field := msg.WhichOneof(od)
	if field == nil {
		return "", errEmptyContainedResource
	}
	return fhirTypeFromOneofField(field), nil
}

func resourceTypeFromR4Message(msg protoreflect.Message) (string, error) {
	if defaultR4Registry.isContainedResource(msg.Descriptor().FullName()) {
		return resourceTypeFromR4ContainedResource(msg.Interface().(*rpb.ContainedResource))
	}
	if rt, ok := defaultR4Registry.fhirTypeForMessage(msg.Descriptor().FullName()); ok {
		return rt, nil
	}
	return "", errUnsupportedProto
}

func isGoogleR4ProtoMessage(msg protoreflect.Message) bool {
	name := msg.Descriptor().FullName()
	if defaultR4Registry.isContainedResource(name) {
		return msg.WhichOneof(msg.Descriptor().Oneofs().ByName(r4ContainedResourceOneof)) != nil
	}
	return defaultR4Registry.isKnownMessage(name)
}
