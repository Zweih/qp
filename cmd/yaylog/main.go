package main

import (
	"fmt"
	"os"
	"yaylog/internal/config"
	"yaylog/internal/display"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
)

func startProgressListener(isInteractive bool, progressChan chan pkgdata.ProgressMessage) {
	if isInteractive {
		go func() {
			for msg := range progressChan {
				fmt.Printf("\r%-80s", fmt.Sprintf("[%s] %d%% - %s", msg.Phase, msg.Progress, msg.Description))
			}
		}()
	} else {
		go func() {
			for range progressChan {
				// discard messages in non-tty
			}
		}()
	}
}

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

	progressChan := make(chan pkgdata.ProgressMessage)
	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))
	startProgressListener(isInteractive, progressChan)
	packages = pkgdata.ConcurrentFilters(packages, progressChan, cfg.DateFilter, cfg.SizeFilter, cfg.ExplicitOnly, cfg.DependenciesOnly)
	close(progressChan)

	if isInteractive {
		fmt.Println()
	}

	pkgdata.SortPackages(packages, cfg.SortBy)

	if cfg.Count > 0 && !cfg.AllPackages && len(packages) > cfg.Count {
		cutoffIdx := len(packages) - cfg.Count
		packages = packages[cutoffIdx:]
	}

	display.PrintTable(packages)
}
