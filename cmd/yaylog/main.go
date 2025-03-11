package main

import (
	"fmt"
	"os"
	"sync"
	"yaylog/internal/config"
	out "yaylog/internal/display"
	"yaylog/internal/pipeline"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
)

func main() {
	err := mainWithConfig(&config.CliConfigProvider{})
	if err != nil {
		out.WriteLine(err.Error())
		os.Exit(1)
	}
}

func mainWithConfig(configProvider config.ConfigProvider) error {
	cfg, err := configProvider.GetConfig()
	if err != nil {
		return err
	}

	packages := fetchPackages()

	isInteractive := term.IsTerminal(int(os.Stdout.Fd())) && !cfg.DisableProgress
	var wg sync.WaitGroup

	pipeline := []PipelinePhase{
		{"Calculating reverse dependencies", pkgdata.CalculateReverseDependencies, isInteractive, &wg},
		{"Filtering", pipeline.PreprocessFiltering, isInteractive, &wg},
		{"Sorting", pkgdata.SortPackages, isInteractive, &wg},
	}

	for _, phase := range pipeline {
		packages, err = phase.Run(cfg, packages)
		if err != nil {
			return err
		}

		if len(packages) == 0 {
			out.WriteLine("No packages to display.")
			return nil
		}
	}

	packages = trimPackagesLen(packages, cfg)
	renderOutput(packages, cfg)
	return nil
}

func fetchPackages() []pkgdata.PackageInfo {
	packages, err := pkgdata.FetchPackages()
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: Some packages may be missing due to corrupted package database: %v", err))
	}

	return packages
}

func trimPackagesLen(
	packages []pkgdata.PackageInfo,
	cfg config.Config,
) []pkgdata.PackageInfo {
	if cfg.Count > 0 && !cfg.AllPackages && len(packages) > cfg.Count {
		cutoffIdx := len(packages) - cfg.Count
		packages = packages[cutoffIdx:]
	}

	return packages
}

func renderOutput(packages []pkgdata.PackageInfo, cfg config.Config) {
	if cfg.OutputJson {
		out.RenderJson(packages, cfg.Fields)
		return
	}

	out.RenderTable(packages, cfg.Fields, cfg.ShowFullTimestamp, cfg.HasNoHeaders)
}
