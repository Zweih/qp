package config

import (
	"fmt"
	"strings"
	"yaylog/internal/consts"
)

func parseColumns(
	columnsInput string,
	addColumnsInput string,
	hasAllColumns bool,
) ([]consts.FieldType, error) {
	var specifiedColumnsRaw string
	var columns []consts.FieldType

	switch {
	case columnsInput != "":
		specifiedColumnsRaw = columnsInput
	case addColumnsInput != "":
		specifiedColumnsRaw = addColumnsInput
		fallthrough
	default:
		if hasAllColumns {
			columns = consts.ValidFields
		} else {
			columns = consts.DefaultFields
		}
	}

	if specifiedColumnsRaw != "" {
		specifiedColumns := strings.Split(
			strings.ToLower(strings.TrimSpace(specifiedColumnsRaw)),
			",",
		)

		for _, column := range specifiedColumns {
			fieldType, exists := consts.FieldTypeLookup[strings.TrimSpace(column)]

			if !exists {
				return nil, fmt.Errorf("Error: '%s' is not a valid column", column)
			}

			columns = append(columns, fieldType)
		}
	}

	return columns, nil
}
