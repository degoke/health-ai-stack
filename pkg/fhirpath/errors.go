package fhirpath

import "errors"

var (
	// ErrInvalidInput is returned when the evaluation root is not a supported type.
	ErrInvalidInput = errors.New("fhirpath: invalid evaluation input")

	// ErrExpressionTooLong is returned when an expression exceeds MaxExpressionLen.
	ErrExpressionTooLong = errors.New("fhirpath: expression exceeds maximum length")

	// ErrTooManyResults is returned when evaluation exceeds MaxResultItems.
	ErrTooManyResults = errors.New("fhirpath: result collection exceeds maximum size")

	// ErrNotSingleton is returned when EvalBool or EvalString requires exactly one value.
	ErrNotSingleton = errors.New("fhirpath: result is not a singleton collection")

	// ErrEmptyResult is returned when EvalBool or EvalString requires a non-empty result.
	ErrEmptyResult = errors.New("fhirpath: result collection is empty")

	// ErrTypeMismatch is returned when a singleton result is not coercible to the target type.
	ErrTypeMismatch = errors.New("fhirpath: result type does not match requested coercion")

	// ErrShadowsBuiltin is returned when a custom function name conflicts with a built-in.
	ErrShadowsBuiltin = errors.New("fhirpath: custom function shadows a built-in name")

	// ErrNotSupported is returned for expressions requiring terminology, resolve(), or
	// other capabilities excluded from the MVP contract.
	ErrNotSupported = errors.New("fhirpath: expression requires unsupported capability")
)
