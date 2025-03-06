package config

import "fmt"

func validateFlagCombinations(
	columnsInput string,
	addColumnsInput string,
	hasAllColumns bool,
	explicitOnly bool,
	dependenciesOnly bool,
) error {
	if columnsInput != "" {
		if addColumnsInput != "" {
			return fmt.Errorf("Error: Cannot use --columns and --add-columns together. Use --columns to fully define the output columns")
		}

		if hasAllColumns {
			return fmt.Errorf("Error: Cannot use --columns and --add-columns together. Use --columns to fully define the output columns")
		}
	}

	if explicitOnly && dependenciesOnly {
		return fmt.Errorf("Error: cannot use --explicit and --dependencies at the same time")
	}

	return nil
}
