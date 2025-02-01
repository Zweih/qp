package main

import (
	"fmt"
	"os"
	"strings"
	"yaylog/internal/config"
	"yaylog/internal/display"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
)

func main() {
	cfg := parseConfig()
	packages := fetchPackages()

	validateConfig(cfg)

	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))

	pipeline := []PipelinePhase{
		{"Filtering", applyFilters, isInteractive},
		{"Sorting", sortPackages, isInteractive},
	}

	for _, phase := range pipeline {
		packages = phase.Run(cfg, packages)

		if len(packages) == 0 {
			fmt.Println("\nNo packages remain after filtering.")
			return
		}
	}

	if cfg.Count > 0 && !cfg.AllPackages && len(packages) > cfg.Count {
		cutoffIdx := len(packages) - cfg.Count
		packages = packages[cutoffIdx:]
	}

	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
	display.PrintTable(packages)
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
		fmt.Printf("Error fetching packages: %v\n", err)
		os.Exit(1)
	}

	return packages
}

func validateConfig(cfg config.Config) {
	if cfg.ExplicitOnly && cfg.DependenciesOnly {
		fmt.Println("Error: Cannot use both --explicit and --dependencies at the same time.")
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
	return pkgdata.SortPackages(packages, cfg.SortBy, reportProgress)
}
