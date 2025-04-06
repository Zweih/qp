package filtering

import (
	"fmt"
	"qp/internal/config"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"sort"
	"strconv"
	"strings"
)

type (
	PkgInfo         = pkgdata.PkgInfo
	FilterCondition = pkgdata.FilterCondition
)

func QueriesToConditions(queries config.FieldQueries) (
	[]*FilterCondition,
	error,
) {
	conditions := make([]*FilterCondition, 0, len(queries))

	for field, subfields := range queries {
		var condition *FilterCondition
		var err error

		switch field {
		case consts.FieldDate:
			condition, err = parseDateCondition(subfields)

		case consts.FieldSize:
			condition, err = parseSizeCondition(subfields)

		case consts.FieldName, consts.FieldArch, consts.FieldLicense:
			condition, err = parseStringCondition(field, subfields)

		case consts.FieldRequiredBy, consts.FieldDepends,
			consts.FieldProvides, consts.FieldConflicts:
			condition, err = parseRelationCondition(field, subfields)

		case consts.FieldReason:
			condition, err = parseReasonCondition(subfields)

		default:
			err = fmt.Errorf("unsupported filter type: %s", consts.FieldNameLookup[field])
		}

		if err != nil {
			return []*FilterCondition{}, err
		}

		conditions = append(conditions, condition)
	}

	// sort filters in order of efficiency
	sort.Slice(conditions, func(i int, j int) bool {
		return conditions[i].FieldType < conditions[j].FieldType
	})

	return conditions, nil
}

func parseRelationCondition(
	field consts.FieldType,
	subfields config.SubfieldQueries,
) (*FilterCondition, error) {
	targetString, hasTarget := subfields[consts.SubfieldTarget]
	if !hasTarget {
		return nil, fmt.Errorf("relation query %s requires target subfield", consts.FieldNameLookup[field])
	}

	targetNames := strings.Split(targetString, ",")
	depthString, hasDepth := subfields[consts.SubfieldDepth]
	var depth int64 = 1
	var err error

	if hasDepth {
		depth, err = strconv.ParseInt(depthString, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid depth value: %s", depthString)
		}
	}

	return newRelationCondition(field, targetNames, int32(depth))
}

func parseStringCondition(
	field consts.FieldType,
	subfields config.SubfieldQueries,
) (*FilterCondition, error) {
	targetString, exists := subfields[consts.SubfieldTarget]
	if !exists {
		return nil, fmt.Errorf("missing target subfield for field: %s", consts.FieldNameLookup[field])
	}

	targets := strings.Split(targetString, ",")
	return newStringCondition(field, targets)
}

func parseReasonCondition(subfields config.SubfieldQueries) (*FilterCondition, error) {
	installReason, exists := subfields[consts.SubfieldTarget]
	if !exists {
		return nil, fmt.Errorf("missing target subfield for field: reason")
	}
	if installReason != config.ReasonExplicit && installReason != config.ReasonDependency {
		return nil, fmt.Errorf("invalid install reason filter: %s", installReason)
	}

	return newReasonCondition(installReason), nil
}

// TODO: we can merge parseDateFilterCondition and parseSizeFilterCondition into parseRangeFilterCondition
func parseSizeCondition(subfields config.SubfieldQueries) (*FilterCondition, error) {
	targetString, exists := subfields[consts.SubfieldTarget]
	if !exists {
		return nil, fmt.Errorf("missing target subfield for field: date")
	}

	sizeFilter, err := parseSizeFilter(targetString)
	if err != nil {
		return nil, fmt.Errorf("invalid size filter: %v", err)
	}

	if err = validateSizeFilter(sizeFilter); err != nil {
		return nil, err
	}

	return newSizeCondition(sizeFilter), nil
}

func parseDateCondition(subfields config.SubfieldQueries) (*FilterCondition, error) {
	targetString, exists := subfields[consts.SubfieldTarget]
	if !exists {
		return nil, fmt.Errorf("missing target subfield for field: date")
	}

	dateFilter, err := parseDateFilter(targetString)
	if err != nil {
		return nil, fmt.Errorf("invalid date filter: %v", err)
	}

	if err = validateDateFilter(dateFilter); err != nil {
		return nil, err
	}

	return newDateCondition(dateFilter), nil
}
