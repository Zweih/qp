package config

import (
	"fmt"
	"strings"
	"yaylog/internal/consts"
)

func parseColumns(columnsInput string, addColumnsInput string, hasAllColumns bool) ([]string, error) {
	var specifiedColumnsRaw string
	var columns []string

	switch {
	case columnsInput != "":
		specifiedColumnsRaw = columnsInput
	case addColumnsInput != "":
		specifiedColumnsRaw = addColumnsInput
		fallthrough
	case hasAllColumns:
		columns = consts.ValidColumns
	default:
		if hasAllColumns {
			columns = consts.ValidColumns
		} else {
			columns = consts.DefaultColumns
		}
	}

	specifiedColumns, err := validateColumns(strings.ToLower(specifiedColumnsRaw))
	if err != nil {
		return nil, err
	}

	columns = append(columns, specifiedColumns...)

	if len(columns) < 1 {
		return nil, fmt.Errorf("no valid columns selected: use --columns to specify at least one column")
	}

	return columns, nil
}

func validateColumns(columnInput string) ([]string, error) {
	if columnInput == "" {
		return []string{}, nil
	}

	validColumnsSet := map[string]bool{}
	for _, columnName := range consts.ValidColumns {
		validColumnsSet[columnName] = true
	}

	var columns []string

	for _, column := range strings.Split(columnInput, ",") {
		cleanColumn := strings.TrimSpace(column)

		if !validColumnsSet[strings.TrimSpace(column)] {
			return nil, fmt.Errorf("%s is not a valid column", cleanColumn)
		}

		columns = append(columns, cleanColumn)
	}

	return columns, nil
}
