package config

import (
	"fmt"
	"qp/internal/consts"
	"strconv"
	"strings"
)

func parseQueries(queryInputs []string) ([]FieldQuery, error) {
	queries := make([]FieldQuery, 0, len(queryInputs))

	for _, input := range queryInputs {
		query, err := parseQueryInput(input)
		if err != nil {
			return nil, err
		}
		queries = append(queries, query)
	}

	return queries, nil
}

func parseQueryInput(input string) (
	FieldQuery,
	error,
) {
	opStart := -1
	opEnd := -1
	match := consts.MatchFuzzy
	negation := false

	for i := range input {
		if input[i] == ':' {
			return parseExistenceQuery(input, i)
		}

		if input[i] == '=' {

			negation = i >= 1 && input[i-1] == '!'
			opStart = i
			if negation {
				opStart--
			}

			opEnd = i + 1

			if opEnd < len(input) && input[opEnd] == '=' {
				match = consts.MatchStrict
				opEnd++
			}

			break
		}
	}

	if opStart < 0 || opEnd < 0 {
		err := fmt.Errorf("invalid query format: %s. Expected e.g. field=value or field==value", input)
		return FieldQuery{}, err
	}

	field, err := parseField(input[:opStart])
	if err != nil {
		return FieldQuery{}, err
	}

	rawTarget := strings.TrimSpace(input[opEnd:])
	target, depth := extractDepth(rawTarget)

	return FieldQuery{
		Negate: negation,
		Field:  field,
		Match:  match,
		Depth:  depth,
		Target: target,
	}, nil
}

func parseExistenceQuery(input string, colonIdx int) (FieldQuery, error) {
	prefix := input[:colonIdx]
	negation := false

	switch prefix {
	case "has":
	case "not":
		negation = true
	default:
		return FieldQuery{}, fmt.Errorf("invalid existence query: %s", input)
	}

	fieldName, depth := extractDepth(input[colonIdx+1:])
	field, err := parseField(fieldName)
	if err != nil {
		return FieldQuery{}, err
	}

	return FieldQuery{
		IsExistence: true,
		Negate:      negation,
		Field:       field,
		Depth:       depth,
	}, nil
}

func parseField(input string) (consts.FieldType, error) {
	fieldName := strings.TrimSpace(input)
	field, exists := consts.FieldTypeLookup[fieldName]
	if !exists {
		return -1, fmt.Errorf("unknown query field: %s", fieldName)
	}

	return field, nil
}

func extractDepth(input string) (target string, depth int32) {
	parts := strings.SplitN(input, "@", 2)
	target = parts[0]
	depth = 1

	if len(parts) == 2 {
		if d, err := strconv.Atoi(parts[1]); err == nil && d > 0 {
			depth = int32(d)
		}
	}

	return target, depth
}
