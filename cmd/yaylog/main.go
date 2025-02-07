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

	validateConfig(cfg)

	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))
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

	out.PrintTable(packages)
}

func parseConfig() config.Config {
	cfg := config.ParseFlags(os.Args[1:])

	if cfg.ShowHelp {
		config.PrintHelp()
		os.Exit(0)
	}

	return cfg
}

func fetchPackages() []pkgdata.PackageInfo {
	packages, err := pkgdata.FetchPackages()
	if err != nil {
		out.WriteLine(fmt.Sprintf("Error fetching packages: %v", err))
		os.Exit(1)
	}

	return packages
}

func validateConfig(cfg config.Config) {
	if cfg.ExplicitOnly && cfg.DependenciesOnly {
		out.WriteLine("Error: Cannot use both --explicit and --dependencies at the same time.")
		os.Exit(1)
	}
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
		out.WriteLine(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}

	return sortedPackages
}
