package validate

import "regexp"

var fhirIDPattern = regexp.MustCompile(`^[A-Za-z0-9\-\.]{1,64}$`)

func defaultRequiredFields() map[string][]string {
	return map[string][]string{
		"Patient":     {},
		"Observation": {"status"},
		"Bundle":      {"type"},
	}
}

func mergeRequiredFields(custom map[string][]string) map[string][]string {
	merged := defaultRequiredFields()
	for rt, fields := range custom {
		merged[rt] = append([]string(nil), fields...)
	}
	return merged
}
