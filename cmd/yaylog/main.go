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

	pkgs := fetchPackages()

	isInteractive := term.IsTerminal(int(os.Stdout.Fd())) && !cfg.DisableProgress
	var wg sync.WaitGroup

	pipelinePhases := []PipelinePhase{
		{"Calculating reverse dependencies", pkgdata.CalculateReverseDependencies, isInteractive, &wg},
		{"Filtering", pipeline.PreprocessFiltering, isInteractive, &wg},
		{"Sorting", pkgdata.SortPackages, isInteractive, &wg},
	}

	for _, phase := range pipelinePhases {
		pkgs, err = phase.Run(cfg, pkgs)
		if err != nil {
			return err
		}

		if len(pkgs) == 0 {
			out.WriteLine("No packages to display.")
			return nil
		}
	}

	pkgs = trimPackagesLen(pkgs, cfg)
	renderOutput(pkgs, cfg)
	return nil
}

func fetchPackages() []pkgdata.PkgInfo {
	pkgs, err := pkgdata.FetchPackages()
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: Some packages may be missing due to corrupted package database: %v", err))
	}

	return pkgs
}

func trimPackagesLen(
	pkgs []pkgdata.PkgInfo,
	cfg config.Config,
) []pkgdata.PkgInfo {
	if cfg.Count > 0 && !cfg.AllPackages && len(pkgs) > cfg.Count {
		cutoffIdx := len(pkgs) - cfg.Count
		pkgs = pkgs[cutoffIdx:]
	}

	return pkgs
}

func renderOutput(pkgs []pkgdata.PkgInfo, cfg config.Config) {
	if cfg.OutputJson {
		out.RenderJson(pkgs, cfg.Fields)
		return
	}

	out.RenderTable(pkgs, cfg.Fields, cfg.ShowFullTimestamp, cfg.HasNoHeaders)
}
