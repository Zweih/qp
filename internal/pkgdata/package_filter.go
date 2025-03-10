package pkgdata

import (
	"time"
	"yaylog/internal/consts"
)

func NewPackageFilter(fieldType consts.FieldType, packageNames []string) FilterCondition {
	packageFilter := newBaseFilter(fieldType)
	var filterFunc Filter

	switch fieldType {
	case consts.FieldName:
		filterFunc = func(pkg PackageInfo) bool {
			return FilterByNames(pkg, packageNames)
		}
	case consts.FieldRequiredBy:
		filterFunc = func(pkg PackageInfo) bool {
			return FilterByPackages(pkg.RequiredBy, packageNames)
		}
	case consts.FieldDepends:
		filterFunc = func(pkg PackageInfo) bool {
			return FilterByPackages(pkg.Depends, packageNames)
		}
	case consts.FieldProvides:
		filterFunc = func(pkg PackageInfo) bool {
			return FilterByPackages(pkg.Provides, packageNames)
		}
	default:
		filterFunc = func(pkg PackageInfo) bool {
			return false // invalid filter type, always return false
		}
	}

	packageFilter.Filter = filterFunc

	return packageFilter
}

func NewDateFilter(start time.Time, end time.Time, isExact bool) FilterCondition {
	dateFilter := newBaseFilter(consts.FieldDate)

	if isExact {
		dateFilter.Filter = func(pkg PackageInfo) bool {
			return FilterByDate(pkg, start)
		}

		return dateFilter
	}

	adjustedEndDate := end.Add(24 * time.Hour) // ensure full date range
	dateFilter.Filter = func(pkg PackageInfo) bool {
		return FilterByDateRange(pkg, start, adjustedEndDate)
	}

	return dateFilter
}

func NewSizeFilter(startSize int64, endSize int64, isExact bool) FilterCondition {
	sizeFilter := newBaseFilter(consts.FieldSize)

	if isExact {
		sizeFilter.Filter = func(pkg PackageInfo) bool {
			return FilterBySize(pkg, startSize)
		}

		return sizeFilter
	}

	sizeFilter.Filter = func(pkg PackageInfo) bool {
		return FilterBySizeRange(pkg, startSize, endSize)
	}

	return sizeFilter
}

func NewReasonFilter(reason string) FilterCondition {
	reasonFilter := newBaseFilter(consts.FieldReason)
	reasonFilter.Filter = func(pkg PackageInfo) bool {
		return FilterByReason(pkg.Reason, reason)
	}

	return reasonFilter
}
