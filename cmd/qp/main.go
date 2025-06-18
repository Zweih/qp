package main

import (
	"fmt"
	"os"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/origins"
	"qp/internal/pipeline/phase"
	"qp/internal/pkgdata"
	"qp/internal/quipple/compiler"
	"qp/internal/storage"
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

	if cfg.CacheOnly != "" {
		return forkCacheWorker(cfg.CacheOnly)
	}

	if cfg.CacheWorker != "" {
		return rebuildCache(cfg.CacheWorker)
	}

	isInteractive := isInteractive(cfg.DisableProgress)
	cacheBaseDir, err := storage.GetCachePath()
	if err != nil {
		out.WriteLine(fmt.Sprintf("WARNING: failed to set up cache dir: %v", err))
	}

	var pipelines []*phase.Pipeline
	for _, driver := range origins.AvailableDrivers() {
		p := phase.NewPipeline(driver, cfg, cacheBaseDir)
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

	allPkgs = pkgdata.EnrichAcrossOrigins(allPkgs)

	if cfg.QueryExpr != nil {
		allPkgs, err = compiler.RunDAG(cfg.QueryExpr, allPkgs)
		if err != nil {
			return fmt.Errorf("filtering failed: %w", err)
		}
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
	err = renderOutput(allPkgs, cfg)
	if err != nil {
		return err
	}

	return nil
}

func isInteractive(disableProgress bool) bool {
	return term.IsTerminal(int(os.Stdout.Fd())) && !disableProgress
}
