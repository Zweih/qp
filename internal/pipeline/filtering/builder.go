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
		case consts.FieldDate, consts.FieldBuildDate, consts.FieldSize:
			condition, err = parseRangeCondition(query)

		case consts.FieldName, consts.FieldReason, consts.FieldVersion,
			consts.FieldOrigin, consts.FieldArch, consts.FieldLicense,
			consts.FieldDescription, consts.FieldUrl, consts.FieldValidation,
			consts.FieldPkgType, consts.FieldPkgBase, consts.FieldPackager:
			condition, err = parseStringCondition(query)

		case consts.FieldConflicts, consts.FieldReplaces,
			consts.FieldDepends, consts.FieldOptDepends,
			consts.FieldRequiredBy, consts.FieldOptionalFor,
			consts.FieldProvides:
			condition, err = parseRelationCondition(query)

		case consts.FieldGroups:
			condition, err = parseStrArrCondition(query)

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

func parseStrArrCondition(query query.FieldQuery) (*FilterCondition, error) {
	if query.IsExistence {
		return newStrArrExistsCondition(query.Field, query.Negate)
	}

	if query.Target == "" {
		return nil, fmt.Errorf("query %s requires a target", consts.FieldNameLookup[query.Field])
	}

	targets := strings.Split(query.Target, ",")
	return newStrArrCondition(query.Field, targets, query.Match, query.Negate)
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
