// Package fhirpath implements haistack-fhirpath, an in-memory FHIRPath expression
// engine for the haistack runtime.
//
// haistack-fhirpath wraps github.com/verily-src/fhirpath-go behind a stable
// haistack-facing API. Callers compile and evaluate FHIRPath expressions against
// a single FHIR resource at a time without importing the upstream library or
// Google FHIR protobuf paths directly.
//
// # Role in the stack
//
// FHIRPath reads fields inside already-selected resources. FHIR Search finds
// resources. ViewDefinitions build reusable projections. This package is the
// expression layer only:
//
//   - It does not query databases or search indexes.
//   - It does not parse FHIR search query strings.
//   - It does not resolve references across a server (resolve() is rejected).
//   - It does not call terminology services (memberOf() is rejected).
//
// Typical consumers are pkg/search (index field extraction), pkg/view
// (ViewDefinition projections), pkg/ai (typed tool preconditions), and future
// conflict or subscription packages that need path-based field access on envelopes
// already loaded from storage.
//
// # Design principles
//
// Single-resource, in-memory evaluation:
//
//   - Every evaluation path accepts exactly one root resource.
//   - The resource becomes the sole item in the FHIRPath input collection.
//   - Multi-resource collections, database handles, and search queries are rejected.
//
// JSON-first runtime, proto-aware evaluation:
//
//   - types.ResourceEnvelope.JSON remains the canonical stored form across haistack.
//   - When envelope.Proto is populated (via pkg/proto), evaluation uses the typed
//     proto directly without re-parsing JSON.
//   - When envelope.Proto is nil, the engine lazily parses envelope.JSON through
//     Config.ProtoCodec (default: proto.NewGoogleR4Codec()) before evaluation.
//
// Stable public surface:
//
//   - Exported types use haistack names (Engine, Value, Collection).
//   - Upstream Verily and Google FHIR types stay behind Value.Raw() and the
//     internal/verily adapter subpackage.
//   - Compile and CompiledExpression values are safe for concurrent use.
//
// Bounded execution:
//
//   - Expression length and result collection size are capped by configuration.
//   - Context deadlines are honored on every evaluation path.
//   - DefaultTimeout applies only when the caller context has no earlier deadline.
//   - Timeout enforcement is soft (see Timeouts below).
//
// Immutable custom-function registry:
//
//   - Custom functions are bound at engine construction and cannot be changed later.
//   - Custom function names must not shadow built-in FHIRPath functions.
//   - Cache keys are expression strings only because the function table is fixed
//     for the lifetime of the engine.
//
// # Backend
//
// The sole MVP backend is github.com/verily-src/fhirpath-go, which evaluates
// expressions against Google FHIR R4 protobuf messages. The concrete adapter
// lives in pkg/fhirpath/internal/verily and is not importable outside this
// module tree.
//
// Broader upstream FHIRPath support may exist but is not part of the haistack
// MVP contract unless covered by tests in this package. Full normative compliance,
// advanced type checking, terminology-backed functions, expression tracing,
// optimizer work, and hard process isolation are deferred.
//
// # Engine interface
//
//	Compile(expr string) (CompiledExpression, error)
//	Eval(ctx context.Context, expr string, resource any) ([]Value, error)
//	EvalBool(ctx context.Context, expr string, resource any) (bool, error)
//	EvalString(ctx context.Context, expr string, resource any) (string, error)
//
// Compile parses and compiles a FHIRPath expression. Repeated calls with the same
// source string return a cached CompiledExpression. Compile is safe to call
// concurrently on the same Engine.
//
// Eval compiles (or reuses a cached compile) and evaluates the expression,
// returning the full result collection as []Value.
//
// EvalBool and EvalString apply strict singleton coercion (see Coercion rules).
//
// Construct an engine with NewEngine:
//
//	eng, err := fhirpath.NewEngine(fhirpath.Config{})
//	if err != nil {
//	    return err
//	}
//
// # CompiledExpression interface
//
//	Expr() string
//	Eval(ctx context.Context, resource any) ([]Value, error)
//	EvalBool(ctx context.Context, resource any) (bool, error)
//	EvalString(ctx context.Context, resource any) (string, error)
//
// A CompiledExpression holds a parsed expression bound to the custom-function
// table configured on the Engine that created it. It is safe to evaluate
// concurrently from multiple goroutines.
//
// Prefer Compile when the same expression is evaluated many times against
// different resources:
//
//	compiled, err := eng.Compile("Patient.telecom.where(system = 'phone').value")
//	if err != nil {
//	    return err
//	}
//	for _, env := range envelopes {
//	    values, err := compiled.Eval(ctx, env)
//	    // ...
//	}
//
// # Evaluation inputs
//
// Accepted root inputs:
//
//   - *types.ResourceEnvelope — preferred runtime container from pkg/core and stores.
//   - Google FHIR R4 protobuf resources recognized by pkg/proto.IsProtoResource,
//     including individual resource messages (for example *patient_go_proto.Patient)
//     and *ContainedResource values with a populated oneof branch.
//
// Rejected root inputs (return ErrInvalidInput):
//
//   - Raw JSON []byte
//   - map[string]any or other decoded JSON structures
//   - Arbitrary Go structs
//   - nil envelopes without JSON or proto payload
//   - Unsupported or non-R4 proto values
//
// Resource adaptation for envelopes:
//
//   - When envelope.Proto is non-nil and supported, it is unwrapped (ContainedResource
//     values yield the inner resource message) and passed to the backend.
//   - When envelope.Proto is nil, envelope.JSON is parsed through ProtoCodec.
//   - envelope.ResourceType is forwarded to ParseJSONToProto for optional type checks.
//
// # Value and Collection
//
// Evaluation returns a Collection ([]Value). Each Value wraps one FHIRPath result
// item from the backend.
//
// Value helpers:
//
//   - Type() string       — stable type name (Boolean, String, Patient, HumanName, …)
//   - Raw() any           — backend-native item; use only when integrating with proto-aware code
//   - Bool() (bool, error)
//   - String() (string, error)
//   - Float64() (float64, error)
//
// Helper methods return ErrTypeMismatch when the wrapped item is not coercible to
// the requested type. Prefer EvalBool and EvalString on Engine or CompiledExpression
// when a singleton boolean or string is required.
//
// # Coercion rules
//
// EvalBool and EvalString enforce stricter rules than the upstream evaluator:
//
//   - Empty collection        → ErrEmptyResult
//   - More than one item      → ErrNotSingleton
//   - One item, wrong type    → ErrTypeMismatch
//
// Eval returns the full collection without coercion. An empty collection is
// represented as nil []Value.
//
// # Config
//
// Config fields and defaults (zero values replaced at NewEngine):
//
//   - ProtoCodec       — default proto.NewGoogleR4Codec()
//   - Functions        — custom FHIRPath functions; default empty map
//   - FunctionArity    — FHIRPath argument count per custom function name; default 0
//   - CacheSize        — default 256 (DefaultCacheSize)
//   - DefaultTimeout   — default 0 (no implicit deadline beyond caller context)
//   - MaxExpressionLen — default 4096 (DefaultMaxExpressionLen)
//   - MaxResultItems   — default 1024 (DefaultMaxResultItems)
//
// MaxCustomFunctionArity (8) is the maximum supported FunctionArity value.
//
// # Custom functions
//
// Register application-specific FHIRPath functions through Config.Functions.
// Each entry maps a FHIRPath function name to a Go callback:
//
//	type Function func(input Collection, args ...Collection) (Collection, error)
//
// The input collection is the current FHIRPath focus collection. Each entry in
// args is the evaluated result of one FHIRPath argument expression, wrapped as a
// haistack Collection (one collection per argument).
//
// Set Config.FunctionArity[name] to the number of FHIRPath arguments the function
// accepts. Omitted names default to arity 0 (niladic calls such as customFn()).
// Each FHIRPath argument must evaluate to a singleton upstream; the adapter wraps
// that singleton as a one-item Collection before invoking the callback.
//
// Example — niladic custom function:
//
//	eng, err := fhirpath.NewEngine(fhirpath.Config{
//	    Functions: map[string]fhirpath.Function{
//	        "alwaysTrue": func(_ fhirpath.Collection, _ ...fhirpath.Collection) (fhirpath.Collection, error) {
//	            return fhirpath.Collection{fhirpath.NewValue(system.Boolean(true))}, nil
//	        },
//	    },
//	})
//
// Example — unary custom function:
//
//	eng, err := fhirpath.NewEngine(fhirpath.Config{
//	    Functions: map[string]fhirpath.Function{
//	        "echoPrefix": func(_ fhirpath.Collection, args ...fhirpath.Collection) (fhirpath.Collection, error) {
//	            s, err := args[0][0].String()
//	            if err != nil {
//	                return nil, err
//	            }
//	            return fhirpath.Collection{fhirpath.NewValue(system.String("prefix-" + s))}, nil
//	        },
//	    },
//	    FunctionArity: map[string]int{"echoPrefix": 1},
//	})
//	// Expression: echoPrefix('hello')  →  "prefix-hello"
//
// Custom function names that conflict with built-in FHIRPath functions cause
// NewEngine to return ErrShadowsBuiltin.
//
// # Timeouts
//
// Evaluation honors context cancellation and deadlines on every path. When
// Config.DefaultTimeout is greater than zero and the caller context has no
// existing deadline, NewEngine-derived evaluation wraps the call in a timeout
// context of that duration.
//
// Because the upstream evaluator has no native cancellation hook, timeout
// enforcement is soft: the wrapper returns context.DeadlineExceeded when the
// deadline is reached but does not force-stop upstream computation already
// running in a background goroutine.
//
// # Errors
//
// Stable sentinel errors:
//
//   - ErrInvalidInput        — unsupported evaluation root type
//   - ErrExpressionTooLong   — expression exceeds MaxExpressionLen
//   - ErrTooManyResults      — result collection exceeds MaxResultItems
//   - ErrEmptyResult         — EvalBool/EvalString on empty collection
//   - ErrNotSingleton        — EvalBool/EvalString on multi-item collection
//   - ErrTypeMismatch        — EvalBool/EvalString or Value helper coercion failure
//   - ErrShadowsBuiltin      — custom function name conflicts with a built-in
//   - ErrNotSupported        — expression requires resolve(), terminology, or other
//     capabilities excluded from the MVP contract
//
// Compile and evaluation may also return upstream parse or evaluation errors,
// context.Canceled, and context.DeadlineExceeded.
//
// # MVP expression subset
//
// Tests in this package guarantee support for common navigation and collection
// operations, including:
//
//   - Field navigation (Patient.name.given)
//   - Filtering (Patient.telecom.where(system = 'phone').value)
//   - Type filtering (Observation.value.ofType(Quantity).value)
//   - Existence, subsetting, and aggregates (exists(), where(), first(), count())
//   - Comparisons and boolean operators
//
// Expressions that call resolve() or experimental terminology functions such as
// memberOf() return ErrNotSupported.
//
// # Integration with other packages
//
// pkg/types  — ResourceEnvelope is the primary evaluation input; JSON is canonical.
// pkg/proto  — GoogleR4Codec parses envelope JSON when Proto is nil; IsProtoResource
// gates direct proto inputs.
// pkg/search — future Indexer implementations may use fhirpath to extract index fields.
// pkg/view   — future ViewDefinition execution may project envelope fields via FHIRPath.
// pkg/ai     — future typed tools may evaluate preconditions on loaded resources.
// pkg/core   — does not depend on fhirpath directly; resources reach fhirpath after load.
//
// pkg/fhirpath must not import pkg/store, pkg/sqlite, pkg/postgres, or other query
// backends.
//
// # Typical flows
//
// Evaluate a path on a stored envelope:
//
//	eng, _ := fhirpath.NewEngine(fhirpath.Config{})
//	given, err := eng.EvalString(ctx, "Patient.name.given", envelope)
//
// Extract multiple values for indexing:
//
//	values, err := eng.Eval(ctx, "Patient.telecom.where(system = 'phone').value", envelope)
//	for _, v := range values {
//	    phone, _ := v.String()
//	    // index phone
//	}
//
// Reuse a compiled expression in a hot loop:
//
//	compiled, _ := eng.Compile("Observation.value.ofType(Quantity).value")
//	for _, obs := range observations {
//	    values, err := compiled.Eval(ctx, obs)
//	    // ...
//	}
//
// # File layout
//
//   - doc.go           — package documentation (this file)
//   - api.go           — Engine, CompiledExpression, Config, Function, Collection
//   - errors.go        — stable sentinel errors
//   - value.go         — Value type and coercion helpers
//   - engine.go        — NewEngine and evaluation orchestration
//   - adapt.go         — custom-function adapter and error mapping
//   - cache.go         — thread-safe compile cache
//   - fhirpath_test.go — MVP contract tests
//   - internal/verily/ — Verily backend adapter (compile, evaluate, resource input)
//
// # Out of scope (MVP)
//
// FHIR search execution, database querying, multi-resource root collections, raw
// JSON evaluation inputs, external constant variables (%name), terminology service
// integration, reference resolution (resolve()), expression tracing, R5 proto
// providers, hard process isolation, and CLI commands (future cmd/haistack fhirpath
// eval). Full FHIRPath normative compliance beyond the tested subset is not
// guaranteed.
package fhirpath
