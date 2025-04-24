package main

import (
	"fmt"
	"os"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/origins"
	"qp/internal/pipeline/phase"
	"qp/internal/pkgdata"
	"qp/internal/syntax"
	"sync"

	"github.com/spf13/pflag"
	"golang.org/x/term"
)

func main() {
	err := mainWithConfig(&config.CliConfigProvider{})
	if err != nil {
		out.WriteLine(fmt.Sprintf("ERROR: %v\n", err.Error()))
		pflag.PrintDefaults()
		os.Exit(1)
	}
}

func mainWithConfig(configProvider config.ConfigProvider) error {
	cfg, err := configProvider.GetConfig()
	if err != nil {
		return err
	}

	isInteractive := isInteractive(cfg.DisableProgress)

	cacheBasePath, err := pkgdata.GetCacheBasePath()
	if err != nil {
		out.WriteLine(fmt.Sprintf("WARNING: failed to set up cache dir: %v", err))
	}

	var pipelines []*phase.Pipeline
	for _, driver := range origins.AvailableDrivers() {
		p := phase.NewPipeline(driver, cfg, isInteractive, cacheBasePath)
		pipelines = append(pipelines, p)
	}

	if len(pipelines) == 0 {
		return fmt.Errorf("no supported package origins detected")
	}

	var wg sync.WaitGroup
	results := make(chan []*pkgdata.PkgInfo, len(pipelines))

	for _, p := range pipelines {
		wg.Add(1)
		go func(p *phase.Pipeline) {
			defer wg.Done()
			pkgs, err := p.Run()
			if err != nil {
				out.WriteLine(fmt.Sprintf("WARNING: [%s] pipeline failed: %v", p.Origin.Name(), err))
				return
			}
			results <- pkgs
		}(p)
	}

	wg.Wait()
	close(results)

	var allPkgs []*pkgdata.PkgInfo
	for pkgs := range results {
		allPkgs = append(allPkgs, pkgs...)
	}

	if len(allPkgs) == 0 {
		if isInteractive {
			out.WriteLine("No packages to display.")
		}

		return nil
	}

	allPkgs, err = globalPackageSort(allPkgs, cfg)
	if err != nil {
		return err
	}

	allPkgs = trimPackagesLen(allPkgs, cfg)
	renderOutput(allPkgs, cfg)

	return nil
}

func globalPackageSort(
	allPkgs []*pkgdata.PkgInfo,
	cfg *config.Config,
) ([]*pkgdata.PkgInfo, error) {
	comparator, err := pkgdata.GetComparator(cfg.SortOption.Field, cfg.SortOption.Asc)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	if len(allPkgs) >= pkgdata.ConcurrentSortThreshold {
		return pkgdata.SortConcurrently(allPkgs, comparator, "", nil), nil
	}

	return pkgdata.SortNormally(allPkgs, comparator, "", nil), nil
}

func trimPackagesLen(
	pkgs []*pkgdata.PkgInfo,
	cfg *config.Config,
) []*pkgdata.PkgInfo {
	if cfg.Limit < 1 || len(pkgs) <= cfg.Limit {
		return pkgs
	}

	switch cfg.LimitMode {
	case syntax.LimitEnd:
		return pkgs[:cfg.Limit]
	case syntax.LimitMid:
		start := (len(pkgs) - cfg.Limit) / 2
		end := start + cfg.Limit
		return pkgs[start:end]
	case syntax.LimitStart:
		fallthrough
	default:
		cutoffIdx := len(pkgs) - cfg.Limit
		return pkgs[cutoffIdx:]
	}
}

func renderOutput(pkgs []*pkgdata.PkgInfo, cfg *config.Config) {
	if cfg.OutputJSON {
		out.RenderJSON(pkgs, cfg.Fields)
		return
	}

	out.RenderTable(pkgs, cfg.Fields, cfg.ShowFullTimestamp, cfg.HasNoHeaders)
}

func isInteractive(disableProgress bool) bool {
	return term.IsTerminal(int(os.Stdout.Fd())) && !disableProgress
}
