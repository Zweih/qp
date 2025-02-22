package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"yaylog/internal/config"
	out "yaylog/internal/display"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
)

func main() {
	cfg := parseConfig()
	packages := fetchPackages()

	err := validateConfig(cfg)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Configuration error: %v", err))
	}

	isInteractive := term.IsTerminal(int(os.Stdout.Fd())) && !cfg.DisableProgress
	var wg sync.WaitGroup

	pipeline := []PipelinePhase{
		{"Filtering", applyFilters, isInteractive, &wg},
		{"Sorting", sortPackages, isInteractive, &wg},
	}

	for _, phase := range pipeline {
		packages = phase.Run(cfg, packages)

		if len(packages) == 0 {
			out.WriteLine("No packages to display.")
			return
		}
	}

	if cfg.Count > 0 && !cfg.AllPackages && len(packages) > cfg.Count {
		cutoffIdx := len(packages) - cfg.Count
		packages = packages[cutoffIdx:]
	}

	out.PrintTable(packages, cfg.ShowFullTimestamp, cfg.OptionalColumns)
}

func parseConfig() config.Config {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		out.WriteLine(fmt.Sprintf("Error parsing arguments: %v", err))
		os.Exit(0)
	}

	if cfg.ShowHelp {
		config.PrintHelp()
		os.Exit(0)
	}

	return cfg
}

func fetchPackages() []pkgdata.PackageInfo {
	packages, err := pkgdata.FetchPackages()
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: Some packages may be missing due to corrupted package database: %v", err))
	}

	return packages
}

func validateConfig(cfg config.Config) error {
	if cfg.ExplicitOnly && cfg.DependenciesOnly {
		return fmt.Errorf("Error: Cannot use both --explicit and --dependencies at the same time")
	}

	return nil
}

func applyFilters(
	cfg config.Config,
	packages []pkgdata.PackageInfo,
	reportProgress pkgdata.ProgressReporter,
) []pkgdata.PackageInfo {
	filters := make([]pkgdata.FilterCondition, 0)

	if cfg.ExplicitOnly {
		filters = append(filters, pkgdata.FilterCondition{
			Filter:    pkgdata.FilterExplicit,
			PhaseName: "Filtering explicit only",
		})
	}

	if cfg.DependenciesOnly {
		filters = append(filters, pkgdata.FilterCondition{
			Filter:    pkgdata.FilterDependencies,
			PhaseName: "Filtering dependencies only",
		})
	}

	if !cfg.DateFilter.StartDate.IsZero() || !cfg.DateFilter.EndDate.IsZero() {
		var dateFilter pkgdata.Filter

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

		filters = append(filters, pkgdata.FilterCondition{
			Filter:    dateFilter,
			PhaseName: "Filtering by date",
		})
	}

	if !(cfg.SizeFilter.StartSize == 0 && cfg.SizeFilter.EndSize == 0) {
		var sizeFilter pkgdata.Filter

		fmt.Println(cfg.SizeFilter.IsExactMatch)

		if cfg.SizeFilter.IsExactMatch {
			sizeFilter = func(pkg pkgdata.PackageInfo) bool {
				return pkgdata.FilterBySize(pkg, cfg.SizeFilter.StartSize)
			}
		} else {
			fmt.Println(cfg.SizeFilter.EndSize)
			sizeFilter = func(pkg pkgdata.PackageInfo) bool {
				return pkgdata.FilterBySizeRange(pkg, cfg.SizeFilter.StartSize, cfg.SizeFilter.EndSize)
			}
		}

		filters = append(filters, pkgdata.FilterCondition{
			Filter:    sizeFilter,
			PhaseName: "Filtering by size",
		})
	}

	if len(cfg.NameFilter) > 0 {
		filters = append(filters, pkgdata.FilterCondition{
			Filter: func(pkg pkgdata.PackageInfo) bool {
				return pkgdata.FilterByName(pkg, cfg.NameFilter)
			},
			PhaseName: "Filtering by name",
		})
	}

	fmt.Println(filters)

	return pkgdata.ApplyFilters(packages, filters, reportProgress)
}

func sortPackages(
	cfg config.Config,
	packages []pkgdata.PackageInfo,
	reportProgress pkgdata.ProgressReporter,
) []pkgdata.PackageInfo {
	sortedPackages, err := pkgdata.SortPackages(packages, cfg.SortBy, reportProgress)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Error sorting packages: %v. Displaying unsorted packages.", err))
		return packages
	}

	return sortedPackages
}
