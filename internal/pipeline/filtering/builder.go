package filtering

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"qp/internal/query"
	"sort"
	"strings"
)

type (
	PkgInfo         = pkgdata.PkgInfo
	FilterCondition = pkgdata.FilterCondition
)

func QueriesToConditions(queries []query.FieldQuery) ([]*FilterCondition, error) {
	conditions := make([]*FilterCondition, 0, len(queries))

	for _, query := range queries {
		var condition *FilterCondition
		var err error

		switch query.Field {
		case consts.FieldDate, consts.FieldSize, consts.FieldBuildDate:
			condition, err = parseRangeCondition(query)

		case consts.FieldName, consts.FieldReason, consts.FieldArch,
			consts.FieldLicense, consts.FieldPkgBase, consts.FieldDescription,
			consts.FieldUrl, consts.FieldValidation, consts.FieldPkgType, consts.FieldPackager:
			condition, err = parseStringCondition(query)

		case consts.FieldDepends, consts.FieldOptDepends,
			consts.FieldRequiredBy, consts.FieldOptionalFor,
			consts.FieldProvides, consts.FieldConflicts:
			condition, err = parseRelationCondition(query)

		default:
			err = fmt.Errorf("unsupported filter type: %s", consts.FieldNameLookup[query.Field])
		}

		if err != nil {
			return []*FilterCondition{}, err
		}

		if condition == nil {
			continue
		}

		conditions = append(conditions, condition)
	}

	// sort filters in order of efficiency
	sort.Slice(conditions, func(i int, j int) bool {
		return conditions[i].FieldType < conditions[j].FieldType
	})

	return conditions, nil
}

func parseRelationCondition(query query.FieldQuery) (*FilterCondition, error) {
	if query.IsExistence {
		return newRelationExistsCondition(query.Field, query.Depth, query.Negate)
	}

	if query.Target == "" {
		return nil, fmt.Errorf("relation query %s requires a target", consts.FieldNameLookup[query.Field])
	}

	targets := strings.Split(query.Target, ",")
	return newRelationCondition(query.Field, targets, query.Depth, query.Match, query.Negate)
}

func parseStringCondition(query query.FieldQuery) (*FilterCondition, error) {
	if query.IsExistence {
		return newStringExistsCondition(query.Field, query.Negate)
	}

	if query.Target == "" {
		return nil, fmt.Errorf("query %s requires a target", consts.FieldNameLookup[query.Field])
	}

	targets := strings.Split(query.Target, ",")
	return newStringCondition(query.Field, targets, query.Match, query.Negate)
}

func parseRangeCondition(query query.FieldQuery) (*FilterCondition, error) {
	if query.IsExistence {
		return nil, nil
	}

	var parser func(string) (RangeSelector, error)
	var validator func(RangeSelector) error

	switch query.Field {
	case consts.FieldDate, consts.FieldBuildDate:
		parser = parseDateFilter
		validator = validateDateFilter
	case consts.FieldSize:
		parser = parseSizeFilter
		validator = validateSizeFilter
	default:
		return nil, fmt.Errorf("field %v is not a valid range field", query.Field)
	}

	selector, err := parser(query.Target)
	if err != nil {
		return nil, err
	}

	if err = validator(selector); err != nil {
		return nil, err
	}

	return newRangeCondition(query, selector)
}
