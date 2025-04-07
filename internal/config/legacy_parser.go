package config

import "qp/internal/consts"

func addLegacyQuery(
	queries FieldQueries,
	field consts.FieldType,
	subfield consts.SubfieldType,
	value string,
) {
	if value == "" {
		return
	}

	subfields := queries[field]
	if subfields == nil {
		subfields = make(SubfieldQueries)
	}

	subfields[subfield] = value
	queries[field] = subfields
}

func convertLegacyQueries(
	queries FieldQueries,
	dateFilter string,
	nameFilter string,
	sizeFilter string,
	requiredByFilter string,
	explicitOnly bool,
	dependenciesOnly bool,
) FieldQueries {
	addLegacyQuery(queries, consts.FieldDate, consts.SubfieldTarget, dateFilter)
	addLegacyQuery(queries, consts.FieldName, consts.SubfieldTarget, nameFilter)
	addLegacyQuery(queries, consts.FieldSize, consts.SubfieldTarget, sizeFilter)
	addLegacyQuery(queries, consts.FieldRequiredBy, consts.SubfieldTarget, requiredByFilter)

	if explicitOnly {
		addLegacyQuery(queries, consts.FieldReason, consts.SubfieldTarget, ReasonExplicit)
	}

	if dependenciesOnly {
		addLegacyQuery(queries, consts.FieldReason, consts.SubfieldTarget, ReasonDependency)
	}

	return queries
}
