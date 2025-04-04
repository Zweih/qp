package filtering

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
)

type RangeSelector struct {
	Start   int64
	End     int64
	IsExact bool
}

type ExactFilter func(pkg *PkgInfo, target int64) bool

type RangeFilter func(pkg *PkgInfo, start int64, end int64) bool

type filterableRangeType interface {
	~int64 | ~[]pkgdata.Relation
}

func newCondition(fieldType consts.FieldType) FilterCondition {
	return FilterCondition{
		PhaseName: "Filtering by " + consts.FieldNameLookup[fieldType],
		FieldType: fieldType,
	}
}

func newStringCondition(fieldType consts.FieldType, targets []string) (*FilterCondition, error) {
	conditionFilter := newCondition(fieldType)
	var filterFunc pkgdata.Filter

	for i, target := range targets {
		targets[i] = strings.ToLower(target)
	}

	switch fieldType {
	case consts.FieldName:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByStrings(pkg.Name, targets)
		}
	case consts.FieldArch:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByStrings(pkg.Arch, targets)
		}
	case consts.FieldLicense:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByStrings(pkg.License, targets)
		}
	case consts.FieldDescription:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByStrings(pkg.Description, targets)
		}
	case consts.FieldRequiredBy:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByRelation(pkg.RequiredBy, targets)
		}
	case consts.FieldDepends:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByRelation(pkg.Depends, targets)
		}
	case consts.FieldProvides:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByRelation(pkg.Provides, targets)
		}
	case consts.FieldConflicts:
		filterFunc = func(pkg *PkgInfo) bool {
			return pkgdata.FilterByRelation(pkg.Conflicts, targets)
		}
	default:
		return nil, fmt.Errorf("invalid field for package filter: %s", consts.FieldNameLookup[fieldType])
	}

	conditionFilter.Filter = filterFunc

	return &conditionFilter, nil
}

func newRangeCondition(
	rangeSelector RangeSelector,
	fieldType consts.FieldType,
	exactFunc ExactFilter,
	rangeFunc RangeFilter,
) *FilterCondition {
	condition := newCondition(fieldType)

	if rangeSelector.IsExact {
		condition.Filter = func(pkg *PkgInfo) bool {
			return exactFunc(pkg, rangeSelector.Start)
		}

		return &condition
	}

	condition.Filter = func(pkg *PkgInfo) bool {
		return rangeFunc(pkg, rangeSelector.Start, rangeSelector.End)
	}

	return &condition
}

func newDateCondition(dateFilter RangeSelector) *FilterCondition {
	return newRangeCondition(
		dateFilter,
		consts.FieldDate,
		pkgdata.FilterByDate,
		pkgdata.FilterByDateRange,
	)
}

func newSizeCondition(sizeFilter RangeSelector) *FilterCondition {
	return newRangeCondition(
		sizeFilter,
		consts.FieldSize,
		pkgdata.FilterBySize, // TODO: maybe these two should be in maps that have the Field as a key
		pkgdata.FilterBySizeRange,
	)
}

func newDepthCondition(fieldType consts.FieldType, depthFilter RangeSelector) *FilterCondition {
	return newRangeCondition(
		depthFilter,
		fieldType,
		func(pkg *PkgInfo, target int64) bool {
			return pkgdata.FilterByDepth(pkg.RequiredBy, target)
		},
		func(pkg *PkgInfo, start, end int64) bool {
			return pkgdata.FilterByDepthRange(pkg.RequiredBy, start, end)
		},
	)
}

func newReasonCondition(reason string) *FilterCondition {
	condition := newCondition(consts.FieldReason)
	condition.Filter = func(pkg *PkgInfo) bool {
		return pkgdata.FilterByReason(pkg.Reason, reason)
	}

	return &condition
}
