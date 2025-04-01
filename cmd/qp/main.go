package main

import (
	"os"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pipeline/meta"
	phasekit "qp/internal/pipeline/phase"
	"qp/internal/pkgdata"
	"sync"

	"github.com/spf13/pflag"
	"golang.org/x/term"
)

func main() {
	err := mainWithConfig(&config.CliConfigProvider{})
	if err != nil {
		out.WriteLine(err.Error() + "\n")
		pflag.PrintDefaults()
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
		phasekit.New("Loading cache", phasekit.LoadCacheStep, &wg),
		phasekit.New("Fetching packages", phasekit.FetchStep, &wg),
		phasekit.New("Calculating reverse dependencies", phasekit.ReverseDepStep, &wg),
		phasekit.New("Saving cache", phasekit.SaveCacheStep, &wg),
		phasekit.New("Filtering", phasekit.FilterStep, &wg),
		phasekit.New("Sorting", phasekit.SortStep, &wg),
	}

	var pkgPtrs []*pkgdata.PkgInfo
	for i, phase := range pipelinePhases {
		pkgPtrs, err = phase.Run(cfg, pkgPtrs, pipelineCtx)
		if err != nil {
			return err
		}

		if i > 0 && len(pkgPtrs) == 0 { // only start checking once loading/fetching has completed
			out.WriteLine("No packages to display.")
			return nil
		}
	}

	pkgPtrs = trimPackagesLen(pkgPtrs, cfg)
	renderOutput(pkgPtrs, cfg)

	return nil
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
