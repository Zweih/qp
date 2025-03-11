package config

import (
	"fmt"
)

func validateFlagCombinations(
	columnsInput string,
	addColumnsInput string,
	hasAllColumns bool,
	explicitOnly bool,
	dependenciesOnly bool,
) error {
	if columnsInput != "" && (addColumnsInput != "" || hasAllColumns) {
		return fmt.Errorf("Error: Cannot use --columns and --add-columns or --all-columns together. Use --columns to fully define the output columns")
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
		return fmt.Errorf("Error: Invalid sort option %s", sortBy)
	}

	return nil
}
