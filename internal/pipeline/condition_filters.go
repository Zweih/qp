package pipeline

import (
	"fmt"
	"time"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

type RangeFilter struct {
	Start   int64
	End     int64
	IsExact bool
}

func newBaseCondition(filterType consts.FieldType) FilterCondition {
	return FilterCondition{
		PhaseName: "Filtering by " + string(filterType),
	}
}

func NewPackageCondition(fieldType consts.FieldType, targets []string) (FilterCondition, error) {
	conditionFilter := newBaseCondition(fieldType)
	var filterFunc pkgdata.Filter

	switch fieldType {
	case consts.FieldName:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByStrings(pkg.Name, targets)
		}
	case consts.FieldArch:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByStrings(pkg.Arch, targets)
		}
	case consts.FieldRequiredBy:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByRelation(pkg.RequiredBy, targets)
		}
	case consts.FieldDepends:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByRelation(pkg.Depends, targets)
		}
	case consts.FieldProvides:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByRelation(pkg.Provides, targets)
		}
	case consts.FieldConflicts:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByRelation(pkg.Conflicts, targets)
		}
	default:
		return FilterCondition{}, fmt.Errorf("invalid field for package filter: %s", fieldType)
	}

	conditionFilter.Filter = filterFunc

	return conditionFilter, nil
}

func NewDateCondition(dateFilter RangeFilter) FilterCondition {
	start, end, isExact := dateFilter.Start, dateFilter.End, dateFilter.IsExact
	condition := newBaseCondition(consts.FieldDate)

	if isExact {
		condition.Filter = func(pkg PackageInfo) bool {
			return pkgdata.FilterByDate(pkg, start)
		}

		return condition
	}

	adjustedEnd := end + int64(time.Hour*24/time.Second) // ensure full date range
	condition.Filter = func(pkg PackageInfo) bool {
		return pkgdata.FilterByDateRange(pkg, start, adjustedEnd)
	}

	return condition
}

func NewSizeCondition(sizeFilter RangeFilter) FilterCondition {
	start, end, isExact := sizeFilter.Start, sizeFilter.End, sizeFilter.IsExact
	condition := newBaseCondition(consts.FieldSize)

	if isExact {
		condition.Filter = func(pkg PackageInfo) bool {
			return pkgdata.FilterBySize(pkg, start)
		}

		return condition
	}

	condition.Filter = func(pkg PackageInfo) bool {
		return pkgdata.FilterBySizeRange(pkg, start, end)
	}

	return condition
}

func NewReasonCondition(reason string) FilterCondition {
	condition := newBaseCondition(consts.FieldReason)
	condition.Filter = func(pkg PackageInfo) bool {
		return pkgdata.FilterByReason(pkg.Reason, reason)
	}

	return condition
}
