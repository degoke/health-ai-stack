package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/degoke/health-ai-stack/pkg/proto"
	"github.com/degoke/health-ai-stack/pkg/types"
)

type builtinEngine struct {
	protoCodec         proto.ProtoCodec
	knownResourceTypes map[string]struct{}
	installedTypes     ResourceTypeRegistry
	requiredFields     map[string][]string
}

// NewEngine returns the built-in haistack-validate engine.
func NewEngine(cfg Config) (Engine, error) {
	if cfg.ProtoCodec == nil {
		cfg.ProtoCodec = proto.NewGoogleR4Codec()
	}
	if cfg.KnownResourceTypes == nil {
		cfg.KnownResourceTypes = proto.KnownR4ResourceTypes()
	}
	cfg.RequiredFields = mergeRequiredFields(cfg.RequiredFields)
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &builtinEngine{
		protoCodec:         cfg.ProtoCodec,
		knownResourceTypes: cfg.KnownResourceTypes,
		installedTypes:     cfg.InstalledTypes,
		requiredFields:     cfg.RequiredFields,
	}, nil
}

func (e *builtinEngine) Validate(ctx context.Context, res *types.ResourceEnvelope, opts ValidateOptions) (*ValidationResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if res == nil {
		return invalidResult(issue("invalid", "resource is nil", nil)), nil
	}
	if len(res.JSON) == 0 {
		return invalidResult(issue("invalid", "resource JSON is empty", nil)), nil
	}

	var issues []ValidationIssue

	obj, jsonErr := decodeJSONObject(res.JSON)
	if jsonErr != nil {
		issues = append(issues, issue("invalid-json", jsonErr.Error(), nil))
		return &ValidationResult{Valid: false, Issues: issues}, nil
	}

	resourceType, rtErr := types.GetResourceType(res.JSON)
	if rtErr != nil {
		issues = append(issues, issue("missing-resource-type", rtErr.Error(), []string{"Resource.resourceType"}))
	} else {
		if _, known := e.knownResourceTypes[resourceType]; !known {
			issues = append(issues, issue(
				"unknown-resource-type",
				fmt.Sprintf("unknown FHIR resource type %q", resourceType),
				[]string{"Resource.resourceType"},
			))
		}

		registry := opts.ResourceTypeRegistry
		if registry == nil {
			registry = e.installedTypes
		}
		if registry != nil {
			if _, known := e.knownResourceTypes[resourceType]; known && !registry.IsInstalled(resourceType) {
				issues = append(issues, issue(
					"resource-type-not-installed",
					fmt.Sprintf("resource type %q is not installed", resourceType),
					[]string{"Resource.resourceType"},
				))
			}
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	id, idErr := types.GetID(res.JSON)
	if idErr != nil {
		issues = append(issues, issue("invalid-id", idErr.Error(), []string{"Resource.id"}))
	} else {
		switch {
		case id == "" && opts.RequireID:
			issues = append(issues, issue("missing-id", "id is required", []string{"Resource.id"}))
		case id != "" && !fhirIDPattern.MatchString(id):
			issues = append(issues, issue(
				"invalid-id",
				fmt.Sprintf("id %q does not match FHIR id syntax", id),
				[]string{"Resource.id"},
			))
		}
	}

	if resourceType != "" {
		if required, ok := e.requiredFields[resourceType]; ok {
			for _, field := range required {
				if !hasTopLevelField(obj, field) {
					issues = append(issues, issue(
						"missing-required-field",
						fmt.Sprintf("required field %q is missing for %s", field, resourceType),
						[]string{resourceType + "." + field},
					))
				}
			}
		}
	}

	if opts.ReferencePolicy == ReferencePolicySyntactic || opts.ReferencePolicy == 0 {
		refs, refErr := types.GetReferences(res.JSON)
		if refErr != nil {
			issues = append(issues, issue("invalid-json", refErr.Error(), nil))
		} else {
			for _, ref := range refs {
				if refIssue := validateReferenceSyntax(ref); refIssue != nil {
					issues = append(issues, *refIssue)
				}
			}
		}
	}

	if shouldRunProtoValidation(resourceType, e.knownResourceTypes, issues) {
		if err := e.validateStructure(ctx, resourceType, res); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			issues = append(issues, issue(
				"structural",
				err.Error(),
				[]string{resourceType},
			))
		}
	}

	if len(issues) > 0 {
		return &ValidationResult{Valid: false, Issues: issues}, nil
	}
	return &ValidationResult{Valid: true}, nil
}

func shouldRunProtoValidation(resourceType string, known map[string]struct{}, issues []ValidationIssue) bool {
	if resourceType == "" {
		return false
	}
	if _, ok := known[resourceType]; !ok {
		return false
	}
	for _, iss := range issues {
		switch iss.Code {
		case "unknown-resource-type", "resource-type-not-installed", "missing-resource-type", "invalid-json":
			return false
		}
	}
	return true
}

func decodeJSONObject(data []byte) (map[string]interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	obj, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected JSON object")
	}
	return obj, nil
}

func hasTopLevelField(obj map[string]interface{}, field string) bool {
	value, ok := obj[field]
	if !ok {
		return false
	}
	if value == nil {
		return false
	}
	return true
}

func issue(code, diagnostics string, expression []string) ValidationIssue {
	return ValidationIssue{
		Severity:    "error",
		Code:        code,
		Diagnostics: diagnostics,
		Expression:  expression,
	}
}

func invalidResult(iss ValidationIssue) *ValidationResult {
	return &ValidationResult{Valid: false, Issues: []ValidationIssue{iss}}
}

func joinIssueDiagnostics(issues []ValidationIssue) string {
	parts := make([]string, 0, len(issues))
	for _, iss := range issues {
		if strings.TrimSpace(iss.Diagnostics) != "" {
			parts = append(parts, iss.Diagnostics)
		}
	}
	return strings.Join(parts, "; ")
}
