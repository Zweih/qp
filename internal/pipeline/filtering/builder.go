package filtering

import (
	"fmt"
	"qp/internal/config"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"sort"
	"strings"
)

type (
	PkgInfo         = pkgdata.PkgInfo
	FilterCondition = pkgdata.FilterCondition
)

func QueriesToConditions(filterQueries map[consts.FilterKey]string) (
	[]*FilterCondition,
	error,
) {
	conditions := make([]*FilterCondition, 0, len(filterQueries))

	for filterKey, value := range filterQueries {
		var condition *FilterCondition
		var err error

		switch filterKey.Subfield {
		case consts.SubfieldNone:
			condition, err = parseBaseCondition(filterKey.Field, value)
		case consts.SubfieldName:
			condition, err = parseStringFilterCondition(filterKey.Field, value)
		case consts.SubfieldDepth:
			condition, err = parseDepthFilterCondition(filterKey.Field, value)
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

func parseDepthFilterCondition(fieldType consts.FieldType, value string) (*FilterCondition, error) {
	depthFilter, err := parseIntRangeFilter(value) // TODO: make a generic int range filter
	if err != nil {
		return nil, fmt.Errorf("invalid depth filter: %v", err)
	}

	return newDepthCondition(fieldType, depthFilter), nil
}

func parseBaseCondition(fieldType consts.FieldType, value string) (*FilterCondition, error) {
	switch fieldType {
	case consts.FieldDate:
		return parseDateFilterCondition(value)
	case consts.FieldSize:
		return parseSizeFilterCondition(value)
	case consts.FieldReason:
		return parseReasonFilterCondition(value)
	case consts.FieldName, consts.FieldRequiredBy, consts.FieldDepends,
		consts.FieldProvides, consts.FieldConflicts, consts.FieldArch,
		consts.FieldLicense, consts.FieldDescription:
		return parseStringFilterCondition(fieldType, value)
	default:
		return nil, fmt.Errorf("unsupported base filter: %s", consts.FieldNameLookup[fieldType])
	}
}

func parseStringFilterCondition(
	fieldType consts.FieldType,
	targetListInput string,
) (*FilterCondition, error) {
	targetList := strings.Split(targetListInput, ",")
	return newStringCondition(fieldType, targetList)
}

func parseReasonFilterCondition(installReason string) (*FilterCondition, error) {
	if installReason != config.ReasonExplicit && installReason != config.ReasonDependency {
		return nil, fmt.Errorf("invalid install reason filter: %s", installReason)
	}

	return newReasonCondition(installReason), nil
}

// TODO: we can merge parseDateFilterCondition and parseSizeFilterCondition into parseRangeFilterCondition
func parseDateFilterCondition(value string) (*FilterCondition, error) {
	dateFilter, err := parseDateFilter(value)
	if err != nil {
		return nil, fmt.Errorf("invalid date filter: %v", err)
	}

	if err = validateDateFilter(dateFilter); err != nil {
		return nil, err
	}

	return newDateCondition(dateFilter), nil
}

func parseSizeFilterCondition(value string) (*FilterCondition, error) {
	sizeFilter, err := parseSizeFilter(value)
	if err != nil {
		return nil, fmt.Errorf("invalid size filter: %v", err)
	}

	if err = validateSizeFilter(sizeFilter); err != nil {
		return nil, err
	}

	return newSizeCondition(sizeFilter), nil
}
