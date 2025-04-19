package syntax

import (
	"fmt"
	"qp/internal/consts"
	"strconv"
	"strings"
)

type FieldQuery struct {
	IsExistence bool
	Negate      bool
	Field       consts.FieldType
	Match       consts.MatchType
	Depth       int32
	Target      string
}

// TODO: stopgap before logic is fleshed out
func ParseQueriesBlock(tokens []string) ([]FieldQuery, error) {
	var queries []FieldQuery
	expectQuery := true

	for i, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		if strings.ToLower(token) == "and" {
			if expectQuery {
				return nil, fmt.Errorf("unexpected 'and' at position %d", i)
			}

			expectQuery = true
			continue
		}

		if !expectQuery {
			return nil, fmt.Errorf("missing 'and' between queries near: %q", token)
		}

		query, err := parseQueryInput(token)
		if err != nil {
			return nil, err
		}

		queries = append(queries, query)
		expectQuery = false
	}

	if expectQuery && len(tokens) > 0 {
		return nil, fmt.Errorf("trailing 'and' with no following condition")
	}

	return queries, nil
}

func ParseQueries(queryInputs []string) ([]FieldQuery, error) {
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
	var depth int32 = 1

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

	target := strings.TrimSpace(input[opEnd:])

	if _, exists := consts.RelationFields[field]; exists {
		target, depth = extractDepth(target)
	}

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
	case "no", "not": // TODO: "not" is legacy, to deprecate
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
