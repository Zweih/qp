package config

import (
	"qp/internal/consts"
)

func convertLegacyQueries(
	queries []FieldQuery,
	dateFilter string,
	nameFilter string,
	sizeFilter string,
	requiredByFilter string,
	explicitOnly bool,
	dependenciesOnly bool,
) []FieldQuery {
	if dateFilter != "" {
		queries = append(queries, FieldQuery{
			Field:  consts.FieldDate,
			Target: dateFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if nameFilter != "" {
		queries = append(queries, FieldQuery{
			Field:  consts.FieldName,
			Target: nameFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if sizeFilter != "" {
		queries = append(queries, FieldQuery{
			Field:  consts.FieldSize,
			Target: sizeFilter,
			Match:  consts.MatchFuzzy,
		})
	}

	if requiredByFilter != "" {
		queries = append(queries, FieldQuery{
			Field:  consts.FieldRequiredBy,
			Target: requiredByFilter,
			Match:  consts.MatchFuzzy,
			Depth:  1,
		})
	}

	if explicitOnly {
		queries = append(queries, FieldQuery{
			Field:  consts.FieldReason,
			Target: ReasonExplicit,
			Match:  consts.MatchFuzzy,
		})
	}

	if dependenciesOnly {
		queries = append(queries, FieldQuery{
			Field:  consts.FieldReason,
			Target: ReasonDependency,
			Match:  consts.MatchFuzzy,
		})
	}

	return queries
}
