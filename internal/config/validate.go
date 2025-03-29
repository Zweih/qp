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
