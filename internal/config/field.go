package config

import (
	"fmt"
	"qp/internal/consts"
	"strings"
)

func parseSelection(
	fieldInput string,
	addFieldInput string,
	hasAllFields bool,
) ([]consts.FieldType, error) {
	var selectFieldsRaw string
	var fields []consts.FieldType

	switch {
	case fieldInput != "":
		selectFieldsRaw = fieldInput
	case addFieldInput != "":
		selectFieldsRaw = addFieldInput
		fallthrough
	default:
		if hasAllFields {
			fields = consts.ValidFields
		} else {
			fields = consts.DefaultFields
		}
	}

	if selectFieldsRaw != "" {
		cleanSelectFields := strings.ToLower(strings.TrimSpace(selectFieldsRaw))

		for selectField := range strings.SplitSeq(cleanSelectFields, ",") {
			field, exists := consts.FieldTypeLookup[strings.TrimSpace(selectField)]

			if !exists {
				return nil, fmt.Errorf("Error: '%s' is not a valid field for selection", selectField)
			}

			fields = append(fields, field)
		}
	}

	return fields, nil
}
