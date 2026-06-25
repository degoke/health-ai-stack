// Package validate defines FHIR resource validation contracts for haistack-core.
//
// haistack-validate is the validation boundary between raw FHIR envelopes and
// pkg/core write pipelines. It specifies when and how resources are checked before
// lifecycle mutations (id assignment, version assignment, persistence) without owning
// storage, indexing, or transport concerns.
//
// # Design principles
//
// Pre-mutation validation:
//
//   - Validators run in pkg/core before id and version fields are assigned or changed.
//   - A validation failure aborts the entire WriteSession (no partial history, index,
//     or outbox writes).
//   - Validation operates on types.ResourceEnvelope; canonical JSON is in envelope.JSON.
//
// Pluggable implementations:
//
//   - Validator is a single-method interface; applications supply structural validators,
//     profile validators, or no-op implementations.
//   - pkg/core treats a nil Validator as a no-op (writes proceed without validation).
//   - Generated proto validators (pkg/proto) or external libraries can implement Validator
//     without pkg/validate importing them.
//
// Errors propagate as core invalid errors:
//
//   - Validator.ValidateResource errors are wrapped by pkg/core as ErrorKindInvalid.
//   - Callers map the resulting error to types.OperationOutcome via
//     core.OperationOutcomeFromError.
//
// # Validator interface
//
//	ValidateResource(ctx context.Context, resource *types.ResourceEnvelope) error
//
// Implementations should:
//
//   - Return nil when the resource is acceptable for persistence.
//   - Return a descriptive error when validation fails (diagnostics appear in
//     OperationOutcome.issue.diagnostics).
//   - Avoid mutating the envelope unless explicitly documented; core expects validation
//     to be read-only.
//   - Use context for deadlines and cancellation on expensive checks.
//
// # Integration with other packages
//
// pkg/core  — invokes Validator on Create, Update, and transaction bundle entries
// before write-side mutations.
// pkg/types — ResourceEnvelope is the validation input; JSON is the canonical payload.
// pkg/proto — optional Google FHIR R4 proto validators may implement Validator by
// converting or inspecting envelope.Proto when populated.
//
// # Typical flows
//
// Structural validator wired into core:
//
//	svc, _ := core.NewResourceService(core.ResourceServiceConfig{
//	    Resources: db.ResourceStore(),
//	    History:   db.HistoryStore(),
//	    Sessions:  db,
//	    Validator: myValidator{},
//	})
//
// No validation (MVP default):
//
//	svc, _ := core.NewResourceService(core.ResourceServiceConfig{
//	    Resources: db.ResourceStore(),
//	    History:   db.HistoryStore(),
//	    Sessions:  db,
//	    // Validator omitted
//	})
//
// # File layout
//
//   - doc.go      — package documentation (this file)
//   - validate.go — Validator interface
//
// # Out of scope (MVP)
//
// Built-in FHIR structural validation, profile loader, terminology services, invariant
// registries, validation OperationOutcome builders with expression paths, and batch
// validation APIs. This package defines the contract only; implementations live in
// application code or future packages (for example pkg/validate/fhir).
package validate
