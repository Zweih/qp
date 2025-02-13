package main

import (
	"fmt"
	"os"
	"sync"
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
		return fmt.Errorf("Error: Cannot use both --explicit and --dependencies at the same time.")
	}

	return nil
}

func applyFilters(
	cfg config.Config,
	packages []pkgdata.PackageInfo,
	reportProgress pkgdata.ProgressReporter,
) []pkgdata.PackageInfo {
	filters := []pkgdata.FilterCondition{
		{
			Condition: cfg.ExplicitOnly,
			Filter:    pkgdata.FilterExplicit,
			PhaseName: "Filtering explicit packages",
		},
		{
			Condition: cfg.DependenciesOnly,
			Filter:    pkgdata.FilterDependencies,
			PhaseName: "Filtering dependencies",
		},
		{
			Condition: !cfg.DateFilter.IsZero(),
			Filter: func(pkgs []pkgdata.PackageInfo) []pkgdata.PackageInfo {
				return pkgdata.FilterByDate(pkgs, cfg.DateFilter)
			},
			PhaseName: "Filtering by date",
		},
		{
			Condition: cfg.SizeFilter.IsFilter,
			Filter: func(pkgs []pkgdata.PackageInfo) []pkgdata.PackageInfo {
				return pkgdata.FilterBySize(pkgs, cfg.SizeFilter.Operator, cfg.SizeFilter.SizeInBytes)
			},
			PhaseName: "Filtering by size",
		},
	}

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
