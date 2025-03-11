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

func NewPackageCondition(fieldType consts.FieldType, packageNames []string) (FilterCondition, error) {
	packageFilter := newBaseCondition(fieldType)
	var filterFunc pkgdata.Filter

	switch fieldType {
	case consts.FieldName:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByNames(pkg, packageNames)
		}
	case consts.FieldRequiredBy:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByPackages(pkg.RequiredBy, packageNames)
		}
	case consts.FieldDepends:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByPackages(pkg.Depends, packageNames)
		}
	case consts.FieldProvides:
		filterFunc = func(pkg PackageInfo) bool {
			return pkgdata.FilterByPackages(pkg.Provides, packageNames)
		}
	default:
		return FilterCondition{}, fmt.Errorf("invalid field for package filter: %s", fieldType)
	}

	packageFilter.Filter = filterFunc

	return packageFilter, nil
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
