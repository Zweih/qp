package pipeline

import (
	"time"
	"yaylog/internal/config"
	"yaylog/internal/pkgdata"
)

func PreprocessFiltering(
	cfg config.Config,
	packages []pkgdata.PackageInfo,
	reportProgress pkgdata.ProgressReporter,
) []pkgdata.PackageInfo {
	var filters []pkgdata.FilterCondition

	filterConditions := []*pkgdata.FilterCondition{
		getRequiredByFilterCondition(cfg),
		getExplicitFilterCondition(cfg),
		getDependenciesFilterCondition(cfg),
		getDateFilterCondition(cfg),
		getSizeFilterCondition(cfg),
		getNameFilterCondition(cfg),
	}

	for _, condition := range filterConditions {
		if condition != nil {
			filters = append(filters, *condition)
		}
	}

	return pkgdata.FilterPackages(packages, filters, reportProgress)
}

func getRequiredByFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if len(cfg.RequiredByFilter) == 0 {
		return nil
	}
	return &pkgdata.FilterCondition{
		Filter: func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterRequiredBy(pkg, cfg.RequiredByFilter)
		},
		PhaseName: "Filter by required package",
	}
}

func getExplicitFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if !cfg.ExplicitOnly {
		return nil
	}

	return &pkgdata.FilterCondition{
		Filter:    pkgdata.FilterExplicit,
		PhaseName: "Filtering explicit only",
	}
}

func getDependenciesFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if !cfg.DependenciesOnly {
		return nil
	}

	return &pkgdata.FilterCondition{
		Filter:    pkgdata.FilterDependencies,
		PhaseName: "Filtering dependencies only",
	}
}

func getDateFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if cfg.DateFilter.StartDate.IsZero() && cfg.DateFilter.EndDate.IsZero() {
		return nil
	}

	var dateFilter func(pkgdata.PackageInfo) bool

	if cfg.DateFilter.IsExactMatch {
		dateFilter = func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterByDate(pkg, cfg.DateFilter.StartDate)
		}
	} else {
		adjustedEndDate := cfg.DateFilter.EndDate.Add(24 * time.Hour)
		dateFilter = func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterByDateRange(pkg, cfg.DateFilter.StartDate, adjustedEndDate)
		}
	}

	return &pkgdata.FilterCondition{
		Filter:    dateFilter,
		PhaseName: "Filtering by date",
	}
}

func getSizeFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if cfg.SizeFilter.StartSize == 0 && cfg.SizeFilter.EndSize == 0 {
		return nil
	}

	var sizeFilter func(pkgdata.PackageInfo) bool

	if cfg.SizeFilter.IsExactMatch {
		sizeFilter = func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterBySize(pkg, cfg.SizeFilter.StartSize)
		}
	} else {
		sizeFilter = func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterBySizeRange(pkg, cfg.SizeFilter.StartSize, cfg.SizeFilter.EndSize)
		}
	}

	return &pkgdata.FilterCondition{
		Filter:    sizeFilter,
		PhaseName: "Filtering by size",
	}
}

func getNameFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if len(cfg.NameFilter) == 0 {
		return nil
	}

	return &pkgdata.FilterCondition{
		Filter: func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterByName(pkg, cfg.NameFilter)
		},
		PhaseName: "Filtering by name",
	}
}
