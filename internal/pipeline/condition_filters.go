package pipeline

import (
	"fmt"
	"time"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

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

func NewDateCondition(dateFilter DateFilter) FilterCondition {
	startDate, endDate, isExact := dateFilter.StartDate, dateFilter.EndDate, dateFilter.IsExact
	condition := newBaseCondition(consts.FieldDate)

	if isExact {
		condition.Filter = func(pkg PackageInfo) bool {
			return pkgdata.FilterByDate(pkg, startDate)
		}

		return condition
	}

	adjustedEndDate := endDate.Add(24 * time.Hour) // ensure full date range
	condition.Filter = func(pkg PackageInfo) bool {
		return pkgdata.FilterByDateRange(pkg, startDate, adjustedEndDate)
	}

	return condition
}

func NewSizeCondition(sizeFilter SizeFilter) FilterCondition {
	startSize, endSize, isExact := sizeFilter.StartSize, sizeFilter.EndSize, sizeFilter.IsExact
	condition := newBaseCondition(consts.FieldSize)

	if isExact {
		condition.Filter = func(pkg PackageInfo) bool {
			return pkgdata.FilterBySize(pkg, startSize)
		}

		return condition
	}

	condition.Filter = func(pkg PackageInfo) bool {
		return pkgdata.FilterBySizeRange(pkg, startSize, endSize)
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
