package main

import (
	"fmt"
	"os"
	"sync"
	"yaylog/internal/config"
	out "yaylog/internal/display"
	"yaylog/internal/pipeline/filtering"
	"yaylog/internal/pipeline/meta"
	phasekit "yaylog/internal/pipeline/phase"
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

	isInteractive := term.IsTerminal(int(os.Stdout.Fd())) && !cfg.DisableProgress
	pipelineCtx := &meta.PipelineContext{IsInteractive: isInteractive}
	var wg sync.WaitGroup

	pipelinePhases := []phasekit.PipelinePhase{
		phasekit.New("Loading cache", loadCache, &wg),
		phasekit.New("Fetching packages", fetchPackages, &wg),
		phasekit.New("Calculating reverse dependencies", pkgdata.ReverseDependencies, &wg),
		phasekit.New("Saving cache", saveCache, &wg),
		phasekit.New("Filtering", filtering.PreprocessFiltering, &wg),
		phasekit.New("Sorting", pkgdata.SortPackages, &wg),
	}

	var pkgPtrs []*pkgdata.PkgInfo
	for i, phase := range pipelinePhases {
		pkgPtrs, err = phase.Run(cfg, pkgPtrs, pipelineCtx)
		if err != nil {
			return err
		}

		if i > 0 && len(pkgPtrs) == 0 { // only start checking once both fetche
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
	pkgPtrs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*pkgdata.PkgInfo, error) {
	if !pipelineCtx.UsedCache {
		var err error
		pkgPtrs, err = pkgdata.FetchPackages()
		if err != nil {
			out.WriteLine(fmt.Sprintf(
				"Warning: Some packages may be missing due to corrupted package database: %v",
				err,
			))
		}
	}

	return pkgPtrs, nil
}

func loadCache(
	_ config.Config,
	_ []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*pkgdata.PkgInfo, error) {
	pkgPtrs, err := pkgdata.LoadProtoCache()
	if err == nil {
		pipelineCtx.UsedCache = true
	}

	// TODO: use ProgressReporter to report cache status
	return pkgPtrs, nil
}

// TODO: add progress reporting
func saveCache(
	_ config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*pkgdata.PkgInfo, error) {
	if !pipelineCtx.UsedCache {
		// TODO: we can probably save the file concurrently
		err := pkgdata.SaveProtoCache(pkgPtrs)
		if err != nil {
			out.WriteLine(fmt.Sprintf("Error saving cache: %v", err))
		}
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
