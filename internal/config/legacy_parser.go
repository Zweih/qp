package config

import (
	"qp/internal/consts"
	"qp/internal/syntax"
)

func ParseLegacyConfig(
	fieldInput, addFieldInput string,
	allFields bool,
	sortInput string,
	filterInputs []string,
	dateFilter, nameFilter, sizeFilter, requiredByFilter string,
	explicitOnly, dependenciesOnly bool,
) (syntax.ParsedInput, error) {
	fields, err := parseSelection(fieldInput, addFieldInput, allFields)
	if err != nil {
		return syntax.ParsedInput{}, err
	}

	sortOpt, err := syntax.ParseSortOption(sortInput)
	if err != nil {
		return syntax.ParsedInput{}, err
	}

	queries, err := syntax.ParseQueries(filterInputs)
	if err != nil {
		return syntax.ParsedInput{}, err
	}

	queries = convertLegacyQueries(
		queries,
		dateFilter,
		nameFilter,
		sizeFilter,
		requiredByFilter,
		explicitOnly,
		dependenciesOnly,
	)

	return syntax.ParsedInput{
		Fields:       fields,
		FieldQueries: queries,
		SortOption:   sortOpt,
	}, nil
}

func convertLegacyQueries(
	queries []syntax.FieldQuery,
	dateFilter string,
	nameFilter string,
	sizeFilter string,
	requiredByFilter string,
	explicitOnly bool,
	dependenciesOnly bool,
) []syntax.FieldQuery {
	if dateFilter != "" {
		queries = append(queries, syntax.FieldQuery{
			Field:  consts.FieldDate,
			Target: dateFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if nameFilter != "" {
		queries = append(queries, syntax.FieldQuery{
			Field:  consts.FieldName,
			Target: nameFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if sizeFilter != "" {
		queries = append(queries, syntax.FieldQuery{
			Field:  consts.FieldSize,
			Target: sizeFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if requiredByFilter != "" {
		queries = append(queries, syntax.FieldQuery{
			Field:  consts.FieldRequiredBy,
			Target: requiredByFilter,
			Match:  consts.MatchFuzzy,
			Depth:  1,
		})
	}

	if explicitOnly {
		queries = append(queries, syntax.FieldQuery{
			Field:  consts.FieldReason,
			Target: ReasonExplicit,
			Match:  consts.MatchFuzzy,
		})
	}

	if dependenciesOnly {
		queries = append(queries, syntax.FieldQuery{
			Field:  consts.FieldReason,
			Target: ReasonDependency,
			Match:  consts.MatchFuzzy,
		})
	}

	return queries
}
