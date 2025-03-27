package main

import (
	"fmt"
	"os"
	"sync"
	"yaylog/internal/config"
	out "yaylog/internal/display"
	"yaylog/internal/pipeline"
	"yaylog/internal/pipeline/meta"
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

	var pkgPtrs []*pkgdata.PkgInfo

	isInteractive := term.IsTerminal(int(os.Stdout.Fd())) && !cfg.DisableProgress
	pipelineCtx := &meta.PipelineContext{IsInteractive: isInteractive}
	var wg sync.WaitGroup

	pipelinePhases := []PipelinePhase{
		{"Fetching packages", fetchPackages, &wg},
		{"Calculating reverse dependencies", pkgdata.CalculateReverseDependencies, &wg},
		{"Saving cache", saveCache, &wg},
		{"Filtering", pipeline.PreprocessFiltering, &wg},
		{"Sorting", pkgdata.SortPackages, &wg},
	}

	for _, phase := range pipelinePhases {
		pkgPtrs, err = phase.Run(cfg, pkgPtrs, pipelineCtx)
		if err != nil {
			return err
		}

		if len(pkgPtrs) == 0 {
			out.WriteLine("No packages to display.")
			return nil
		}
	}

	pkgPtrs = trimPackagesLen(pkgPtrs, cfg)

	renderOutput(pkgPtrs, cfg)
	return nil
}

// TODO: add progress reporting
func fetchPackages(
	_ config.Config,
	_ []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*pkgdata.PkgInfo, error) {
	// TODO: break these up into separate phases
	pkgPtrs, err := pkgdata.LoadProtoCache()
	if err == nil {
		pipelineCtx.UsedCache = true
		return pkgPtrs, nil
	}

	pkgPtrs, err = pkgdata.FetchPackages()
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: Some packages may be missing due to corrupted package database: %v", err))
	}

	return pkgPtrs, nil
}

// TODO: add progress reporting
func saveCache(
	_ config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
	_ *meta.PipelineContext,
) ([]*pkgdata.PkgInfo, error) {
	// TODO: we can probably save the file concurrently
	err := pkgdata.SaveProtoCache(pkgPtrs)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Error saving cache: %v", err))
	}

	return pkgPtrs, nil
}

func trimPackagesLen(
	pkgPtrs []*pkgdata.PkgInfo,
	cfg config.Config,
) []*pkgdata.PkgInfo {
	if cfg.Count > 0 && !cfg.AllPackages && len(pkgPtrs) > cfg.Count {
		cutoffIdx := len(pkgPtrs) - cfg.Count
		pkgPtrs = pkgPtrs[cutoffIdx:]
	}

	return pkgPtrs
}

func renderOutput(pkgs []*pkgdata.PkgInfo, cfg config.Config) {
	if cfg.OutputJson {
		out.RenderJson(pkgs, cfg.Fields)
		return
	}

	out.RenderTable(pkgs, cfg.Fields, cfg.ShowFullTimestamp, cfg.HasNoHeaders)
}
