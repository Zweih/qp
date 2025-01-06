package main

import (
	"fmt"
	"os"
	"yaylog/internal/config"
	"yaylog/internal/display"
	"yaylog/internal/pkgdata"
)

func main() {
	// parse cli args (excluding the program name itself: os.Args[0])
	cfg := config.ParseFlags(os.Args[1:])

	// on -h or --help: print help and exit
	if cfg.ShowHelp {
		config.PrintHelp()
		return
	}

	// fetch mock package data
	packages, err := pkgdata.FetchPackages()
	if err != nil {
		// if something went wrong in fetching, log it and exit

		fmt.Printf("Error fetching packages: %v\n", err)
		os.Exit(1)
	}

	if cfg.ExplicitOnly && cfg.DependenciesOnly {
		fmt.Println("Error: Cannot use both --explicit and --dependencies at the same time.")
		os.Exit(1)
	}

	if cfg.ExplicitOnly {
		packages = pkgdata.FilterExplicit(packages)
	}

	if cfg.DependenciesOnly {
		packages = pkgdata.FilterDependencies(packages)
	}

	if !cfg.DateFilter.IsZero() {
		packages = pkgdata.FilterByDate(packages, cfg.DateFilter)
	}

	pkgdata.SortPackages(packages, cfg.SortBy)

	if cfg.Count > 0 && len(packages) > cfg.Count {
		packages = packages[:cfg.Count]
	}

	display.PrintTable(packages)
}
