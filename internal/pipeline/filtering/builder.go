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

func QueriesToConditions(filterQueries map[consts.FieldType]string) (
	[]*FilterCondition,
	error,
) {
	conditions := make([]*FilterCondition, 0, len(filterQueries))

	for fieldType, value := range filterQueries {
		var condition *FilterCondition
		var err error

		switch fieldType {
		case consts.FieldDate:
			condition, err = parseDateFilterCondition(value)
		case consts.FieldSize:
			condition, err = parseSizeFilterCondition(value)
		case consts.FieldName, consts.FieldRequiredBy, consts.FieldDepends,
			consts.FieldProvides, consts.FieldConflicts, consts.FieldArch, consts.FieldLicense:
			condition, err = parsePackageFilterCondition(fieldType, value)
		case consts.FieldReason:
			condition, err = parseReasonFilterCondition(value)
		default:
			err = fmt.Errorf("unsupported filter type: %s", consts.FieldNameLookup[fieldType])
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

func parsePackageFilterCondition(
	fieldType consts.FieldType,
	targetListInput string,
) (*FilterCondition, error) {
	targetList := strings.Split(targetListInput, ",")
	return newPackageCondition(fieldType, targetList)
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
