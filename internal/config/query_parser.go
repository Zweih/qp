package config

import (
	"fmt"
	"qp/internal/consts"
	"strings"
)

func parseQueries(queryInputs []string) (FieldQueries, error) {
	queries := make(FieldQueries)

	for _, queryInput := range queryInputs {
		fieldPart, value, err := parseQueryInput(queryInput)
		if err != nil {
			return nil, err
		}

		field, subfield, err := parseFieldPart(fieldPart)
		if err != nil {
			return nil, err
		}

		addQuery(queries, field, subfield, value)
	}

	return queries, nil
}

func parseQueryInput(input string) (string, string, error) {
	queryParts := strings.SplitN(input, "=", 2)
	if len(queryParts) != 2 {
		return "",
			"",
			fmt.Errorf("invalid query format: %s. Must be in form fireld.subfield=value", input)
	}

	return queryParts[0], queryParts[1], nil
}

func parseFieldPart(fieldPart string) (consts.FieldType, consts.SubfieldType, error) {
	fieldParts := strings.SplitN(fieldPart, ".", 2)
	fieldName := fieldParts[0]
	field, exists := consts.FieldTypeLookup[fieldName]
	if !exists {
		return 0, 0, fmt.Errorf("unknown query field: %s", fieldName)
	}

	subfieldName := ""
	if len(fieldParts) == 2 {
		subfieldName = fieldParts[1]
	}

	subfield, exists := consts.SubfieldTypeLookup[subfieldName]
	if !exists {
		return 0, 0, fmt.Errorf("unknown query subfield: %s", subfieldName)
	}

	return field, subfield, nil
}

func addQuery(
	queries FieldQueries,
	field consts.FieldType,
	subfield consts.SubfieldType,
	value string,
) {
	subfields, exists := queries[field]
	if !exists {
		subfields = make(SubfieldQueries)
	}

	subfields[subfield] = value
	queries[field] = subfields
}
