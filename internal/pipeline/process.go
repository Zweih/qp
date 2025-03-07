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

	filterCondition := &pkgdata.FilterCondition{
		PhaseName: "Filtering by date",
	}

	if cfg.DateFilter.IsExactMatch {
		filterCondition.Filter = func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterByDate(pkg, cfg.DateFilter.StartDate)
		}

		return filterCondition
	}

	adjustedEndDate := cfg.DateFilter.EndDate.Add(24 * time.Hour)
	filterCondition.Filter = func(pkg pkgdata.PackageInfo) bool {
		return pkgdata.FilterByDateRange(pkg, cfg.DateFilter.StartDate, adjustedEndDate)
	}

	return filterCondition
}

func getSizeFilterCondition(cfg config.Config) *pkgdata.FilterCondition {
	if cfg.SizeFilter.StartSize == 0 && cfg.SizeFilter.EndSize == 0 {
		return nil
	}

	var sizeFilter func(pkgdata.PackageInfo) bool
	filterCondition := &pkgdata.FilterCondition{
		Filter:    sizeFilter,
		PhaseName: "Filtering by size",
	}

	if cfg.SizeFilter.IsExactMatch {
		filterCondition.Filter = func(pkg pkgdata.PackageInfo) bool {
			return pkgdata.FilterBySize(pkg, cfg.SizeFilter.StartSize)
		}

		return filterCondition
	}

	filterCondition.Filter = func(pkg pkgdata.PackageInfo) bool {
		return pkgdata.FilterBySizeRange(pkg, cfg.SizeFilter.StartSize, cfg.SizeFilter.EndSize)
	}

	return filterCondition
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
