package config

import (
	"fmt"
)

func validateFlagCombinations(
	fieldInput string,
	addFieldInput string,
	hasAllFields bool,
	explicitOnly bool,
	dependenciesOnly bool,
) error {
	if fieldInput != "" && (addFieldInput != "" || hasAllFields) {
		return fmt.Errorf("Error: Cannot use --select/--select-add or --select-al together. Use --select to fully define the output fields")
	}

	if explicitOnly && dependenciesOnly {
		return fmt.Errorf("Error: cannot use --explicit and --dependencies at the same time")
	}

	return nil
}

func validateConfig(cfg Config) error {
	if err := validateSortOption(cfg.SortBy); err != nil {
		return err
	}

	return nil
}

func validateSortOption(sortBy string) error {
	validSortOptions := map[string]bool{
		"date":         true,
		"alphabetical": true,
		"size:desc":    true,
		"size:asc":     true,
	}

	if !validSortOptions[sortBy] {
		return fmt.Errorf("Error: Invalid order option %s", sortBy)
	}

	return nil
}
