package pipeline

import (
	"yaylog/internal/config"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"
)

const (
	FieldExplicit   = "explicit"
	FieldDependency = "dependency"
)

func PreprocessFiltering(
	cfg config.Config,
	packages []pkgdata.PackageInfo,
	reportProgress pkgdata.ProgressReporter,
) []pkgdata.PackageInfo {
	var filters []pkgdata.FilterCondition

	filterConditions := getFiltersFromConfig(cfg)

	for _, condition := range filterConditions {
		if condition != nil {
			filters = append(filters, *condition)
		}
	}

	return pkgdata.FilterPackages(packages, filters, reportProgress)
}

// TODO: remove these if-statements when we consolidate all filters into one flag
func getFiltersFromConfig(cfg config.Config) []*pkgdata.FilterCondition {
	var conditions []*pkgdata.FilterCondition

	if len(cfg.RequiredByFilter) > 0 {
		filter := pkgdata.NewPackageFilter(consts.FieldRequiredBy, []string{cfg.RequiredByFilter})
		conditions = append(conditions, &filter)
	}

	if cfg.ExplicitOnly {
		filter := pkgdata.NewReasonFilter(FieldExplicit)
		conditions = append(conditions, &filter)
	}

	if cfg.DependenciesOnly {
		filter := pkgdata.NewReasonFilter(FieldDependency)
		conditions = append(conditions, &filter)
	}

	if !cfg.DateFilter.StartDate.IsZero() || !cfg.DateFilter.EndDate.IsZero() {
		filter := pkgdata.NewDateFilter(cfg.DateFilter.StartDate, cfg.DateFilter.EndDate, cfg.DateFilter.IsExactMatch)
		conditions = append(conditions, &filter)
	}

	if cfg.SizeFilter.StartSize > 0 || cfg.SizeFilter.EndSize > 0 {
		filter := pkgdata.NewSizeFilter(cfg.SizeFilter.StartSize, cfg.SizeFilter.EndSize, cfg.SizeFilter.IsExactMatch)
		conditions = append(conditions, &filter)
	}

	if len(cfg.NameFilter) > 0 {
		filter := pkgdata.NewPackageFilter(consts.FieldName, []string{cfg.NameFilter})
		conditions = append(conditions, &filter)
	}

	return conditions
}
