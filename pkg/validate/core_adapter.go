package validate

import (
	"context"

	"github.com/degoke/health-ai-stack/pkg/types"
)

type coreValidator struct {
	engine Engine
	opts   ValidateOptions
}

// NewCoreValidator adapts Engine into the pkg/core Validator interface.
func NewCoreValidator(engine Engine, opts ValidateOptions) Validator {
	return &coreValidator{engine: engine, opts: opts}
}

func (v *coreValidator) ValidateResource(ctx context.Context, resource *types.ResourceEnvelope) error {
	result, err := v.engine.Validate(ctx, resource, v.opts)
	if err != nil {
		return err
	}
	if result == nil || result.Valid {
		return nil
	}
	return validationError{issues: result.Issues}
}

type validationError struct {
	issues []ValidationIssue
}

func (e validationError) Error() string {
	return joinIssueDiagnostics(e.issues)
}
