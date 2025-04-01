package config

import (
	"fmt"
	"qp/internal/consts"
	"strings"
)

func parseFields(
	fieldInput string,
	addFieldInput string,
	hasAllFields bool,
) ([]consts.FieldType, error) {
	var specifiedColumnsRaw string
	var fields []consts.FieldType

	switch {
	case fieldInput != "":
		specifiedColumnsRaw = fieldInput
	case addFieldInput != "":
		specifiedColumnsRaw = addFieldInput
		fallthrough
	default:
		if hasAllFields {
			fields = consts.ValidFields
		} else {
			fields = consts.DefaultFields
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

			fields = append(fields, fieldType)
		}
	}

	return fields, nil
}
