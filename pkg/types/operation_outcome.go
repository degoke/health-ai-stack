package types

// OperationOutcome is a hand-written FHIR OperationOutcome resource for error responses.
// Marshal and unmarshal directly with encoding/json; Code is a plain string in MVP.
type OperationOutcome struct {
	ResourceType string           `json:"resourceType"`
	Issue        []OperationIssue `json:"issue,omitempty"`
}

// OperationIssue is one issue entry within an OperationOutcome.
type OperationIssue struct {
	Severity    string   `json:"severity,omitempty"`
	Code        string   `json:"code,omitempty"`
	Diagnostics string   `json:"diagnostics,omitempty"`
	Expression  []string `json:"expression,omitempty"`
}
