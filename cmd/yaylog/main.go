package main

import (
	"fmt"
	"os"
	"yaylog/internal/config"
	"yaylog/internal/display"
	"yaylog/internal/pkgdata"
)

func main() {
	cfg := config.ParseFlags(os.Args[1:])

	// on -h or --help: print help and exit
	if cfg.ShowHelp {
		config.PrintHelp()
		return
	}

	packages, err := pkgdata.FetchPackages()
	if err != nil {
		fmt.Printf("Error fetching packages: %v\n", err)
		os.Exit(1)
	}

	if cfg.ExplicitOnly && cfg.DependenciesOnly {
		fmt.Println("Error: Cannot use both --explicit and --dependencies at the same time.")
		os.Exit(1)
	}

	packages = pkgdata.ConcurrentFilters(packages, cfg.DateFilter, cfg.ExplicitOnly, cfg.DependenciesOnly)
	pkgdata.SortPackages(packages, cfg.SortBy)

	if cfg.Count > 0 && !cfg.AllPackages && len(packages) > cfg.Count {
		cutoffIdx := len(packages) - cfg.Count
		packages = packages[cutoffIdx:]
	}

	display.PrintTable(packages)
}
