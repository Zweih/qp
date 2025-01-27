package main

import (
	"fmt"
	"os"
	"yaylog/internal/config"
	"yaylog/internal/display"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
)

func parseConfig() config.Config {
	cfg := config.ParseFlags(os.Args[1:])

	if cfg.ShowHelp {
		config.PrintHelp()
		os.Exit(0)
	}

	return cfg
}

func main() {
	cfg := parseConfig()

	packages, err := pkgdata.FetchPackages()
	if err != nil {
		fmt.Printf("Error fetching packages: %v\n", err)
		os.Exit(1)
	}

	if cfg.ExplicitOnly && cfg.DependenciesOnly {
		fmt.Println("Error: Cannot use both --explicit and --dependencies at the same time.")
		os.Exit(1)
	}

	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))
	var progressChan chan pkgdata.ProgressMessage

	if isInteractive {
		progressChan = make(chan pkgdata.ProgressMessage)

		go func() {
			for msg := range progressChan {
				fmt.Printf("\r%-80s", fmt.Sprintf("[%s] %d%% - %s", msg.Phase, msg.Progress, msg.Description))
			}
		}()
	}

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

	packages = pkgdata.ApplyFilters(packages, filters, func(current, total int, phase string) {
		if progressChan != nil {
			progressChan <- pkgdata.ProgressMessage{
				Phase:       phase,
				Progress:    (current * 100) / total,
				Description: "Filtering in progress...",
			}
		}
	})

	if isInteractive {
		close(progressChan)
		fmt.Println()
	}

	pkgdata.SortPackages(packages, cfg.SortBy)

	if cfg.Count > 0 && !cfg.AllPackages && len(packages) > cfg.Count {
		cutoffIdx := len(packages) - cfg.Count
		packages = packages[cutoffIdx:]
	}

	display.PrintTable(packages)
}
