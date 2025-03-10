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

	if err := validateDateFilter(cfg.DateFilter); err != nil {
		return err
	}

	if err := validateSizeFilter(cfg.SizeFilter); err != nil {
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

func validateDateFilter(dateFilter DateFilter) error {
	if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
		if dateFilter.StartDate.After(dateFilter.EndDate) {
			return fmt.Errorf("Error invalid date range. The start date cannot be after the end date")
		}
	}

	return nil
}

func validateSizeFilter(sizeFilter SizeFilter) error {
	if sizeFilter.StartSize > 0 && sizeFilter.EndSize > 0 {
		if sizeFilter.StartSize > sizeFilter.EndSize {
			return fmt.Errorf("Error: invalid size range. Start size cannot be greater than the end size")
		}
	}

	return nil
}
