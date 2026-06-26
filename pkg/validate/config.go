package validate

import "fmt"

func validateConfig(cfg *Config) error {
	if cfg.KnownResourceTypes != nil && len(cfg.KnownResourceTypes) == 0 {
		return fmt.Errorf("validate: KnownResourceTypes must not be empty")
	}
	for resourceType, fields := range cfg.RequiredFields {
		if resourceType == "" {
			return fmt.Errorf("validate: RequiredFields contains empty resource type")
		}
		for _, field := range fields {
			if field == "" {
				return fmt.Errorf("validate: RequiredFields for %s contains empty field name", resourceType)
			}
		}
	}
	return nil
}
