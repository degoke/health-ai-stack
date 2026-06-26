package registry

import "time"

// SearchParameterInfo is compiled search parameter metadata for one resource type.
type SearchParameterInfo struct {
	CanonicalURL string `json:"canonicalUrl"`
	Version      string `json:"version"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Expression   string `json:"expression"`
}

// ResourceCapability describes one enabled resource type in the compiled registry.
type ResourceCapability struct {
	ResourceType        string                `json:"resourceType"`
	StructureDefinition *DefinitionRef        `json:"structureDefinition,omitempty"`
	SearchParameters    []SearchParameterInfo `json:"searchParameters,omitempty"`
}

// CapabilitySnapshot is a lightweight capability view for HTTP and runtime consumers.
type CapabilitySnapshot struct {
	FHIRVersion string               `json:"fhirVersion"`
	Resources   []ResourceCapability `json:"resources"`
	CompiledAt  time.Time            `json:"compiledAt"`
}

// DefinitionRef identifies one installed definition by canonical URL and version.
type DefinitionRef struct {
	CanonicalURL string `json:"canonicalUrl"`
	Version      string `json:"version"`
}

// InstallProvenance captures where an installed definition came from.
type InstallProvenance struct {
	PackageName    string
	PackageVersion string
	ModuleName     string
	SourceModule   string
}
