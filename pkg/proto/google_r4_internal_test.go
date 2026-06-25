package proto

import (
	"testing"

	rpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
)

func TestGoogleR4Codec_ProtoToJSON_IndividualResource(t *testing.T) {
	data := []byte(`{"resourceType":"Patient","id":"ind-1","name":[{"text":"Sam"}]}`)
	codec := NewGoogleR4Codec()

	cr, err := codec.ParseJSONToProto("", data)
	if err != nil {
		t.Fatalf("ParseJSONToProto: %v", err)
	}
	patient := cr.(*rpb.ContainedResource).GetPatient()
	if patient == nil {
		t.Fatal("expected patient resource in contained resource")
	}

	out, err := codec.ProtoToJSON("Patient", patient)
	if err != nil {
		t.Fatalf("ProtoToJSON: %v", err)
	}
	if !IsProtoResource(patient) {
		t.Fatal("individual patient proto should be recognized")
	}
	rt, err := ResourceTypeOfProto(patient)
	if err != nil {
		t.Fatalf("ResourceTypeOfProto: %v", err)
	}
	if rt != "Patient" {
		t.Errorf("ResourceTypeOfProto = %q, want Patient", rt)
	}

	envelope, err := codec.ProtoToEnvelope("Patient", patient)
	if err != nil {
		t.Fatalf("ProtoToEnvelope: %v", err)
	}
	if envelope.ID != "ind-1" {
		t.Errorf("ID = %q, want ind-1", envelope.ID)
	}
	if envelope.Proto != patient {
		t.Error("envelope should retain the individual proto value")
	}
	if string(out) != string(envelope.JSON) {
		t.Errorf("ProtoToJSON output should match envelope JSON")
	}
}

func TestWrapR4Resource(t *testing.T) {
	data := []byte(`{"resourceType":"Patient","id":"wrap-1"}`)
	codec := NewGoogleR4Codec()
	cr, err := codec.ParseJSONToProto("", data)
	if err != nil {
		t.Fatalf("ParseJSONToProto: %v", err)
	}
	patient := cr.(*rpb.ContainedResource).GetPatient()

	wrapped, err := wrapR4Resource(patient)
	if err != nil {
		t.Fatalf("wrapR4Resource: %v", err)
	}
	rt, err := ResourceTypeOfProto(wrapped)
	if err != nil {
		t.Fatalf("ResourceTypeOfProto: %v", err)
	}
	if rt != "Patient" {
		t.Errorf("ResourceTypeOfProto = %q, want Patient", rt)
	}
}
