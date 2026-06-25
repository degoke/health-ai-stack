package core

import (
	"errors"

	"github.com/degoke/health-ai-stack/pkg/types"
)

// OperationOutcomeFromError maps a haistack-core error to a FHIR OperationOutcome.
func OperationOutcomeFromError(err error) *types.OperationOutcome {
	if err == nil {
		return nil
	}

	kind := KindOf(err)
	severity := "error"
	code := "exception"
	diagnostics := err.Error()

	switch kind {
	case ErrorKindInvalid:
		code = "invalid"
	case ErrorKindConflict:
		code = "conflict"
	case ErrorKindNotFound:
		code = "not-found"
	case ErrorKindNotSupported:
		code = "not-supported"
	case ErrorKindException:
		code = "exception"
	}

	issue := types.OperationIssue{
		Severity:    severity,
		Code:        code,
		Diagnostics: diagnostics,
	}

	var svcErr *ServiceError
	if errors.As(err, &svcErr) && len(svcErr.Expression) > 0 {
		issue.Expression = append([]string(nil), svcErr.Expression...)
	}

	return &types.OperationOutcome{
		ResourceType: "OperationOutcome",
		Issue:        []types.OperationIssue{issue},
	}
}
