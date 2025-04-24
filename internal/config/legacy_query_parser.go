package config

import (
	"qp/internal/consts"
	"qp/internal/query"
	"qp/internal/syntax"
)

func ParseLegacyConfig(
	fieldInput, addFieldInput string,
	allFields bool,
	sortInput string,
	filterInputs []string,
	dateFilter, nameFilter, sizeFilter, requiredByFilter string,
	explicitOnly, dependenciesOnly, allPackages bool,
	count int,
) (syntax.ParsedInput, error) {
	fields, err := parseSelection(fieldInput, addFieldInput, allFields)
	if err != nil {
		return syntax.ParsedInput{}, err
	}

	sortOpt, err := syntax.ParseSortOption(sortInput)
	if err != nil {
		return syntax.ParsedInput{}, err
	}

	queries, err := query.ParseQueries(filterInputs)
	if err != nil {
		return syntax.ParsedInput{}, err
	}

	if allPackages {
		count = 0
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
		Limit:        count,
	}, nil
}

func convertLegacyQueries(
	queries []query.FieldQuery,
	dateFilter string,
	nameFilter string,
	sizeFilter string,
	requiredByFilter string,
	explicitOnly bool,
	dependenciesOnly bool,
) []query.FieldQuery {
	if dateFilter != "" {
		queries = append(queries, query.FieldQuery{
			Field:  consts.FieldDate,
			Target: dateFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if nameFilter != "" {
		queries = append(queries, query.FieldQuery{
			Field:  consts.FieldName,
			Target: nameFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if sizeFilter != "" {
		queries = append(queries, query.FieldQuery{
			Field:  consts.FieldSize,
			Target: sizeFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if requiredByFilter != "" {
		queries = append(queries, query.FieldQuery{
			Field:  consts.FieldRequiredBy,
			Target: requiredByFilter,
			Match:  consts.MatchFuzzy,
			Depth:  1,
		})
	}

	if explicitOnly {
		queries = append(queries, query.FieldQuery{
			Field:  consts.FieldReason,
			Target: ReasonExplicit,
			Match:  consts.MatchFuzzy,
		})
	}

	if dependenciesOnly {
		queries = append(queries, query.FieldQuery{
			Field:  consts.FieldReason,
			Target: ReasonDependency,
			Match:  consts.MatchFuzzy,
		})
	}

	return queries
}
