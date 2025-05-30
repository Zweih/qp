package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"qp/api/driver"
	"qp/internal/config"
	"qp/internal/consts"
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

	if cfg.CacheOnly != "" {
		return forkCacheWorker(cfg.CacheOnly)
	}

	if cfg.CacheWorker != "" {
		return rebuildCache(cfg.CacheWorker)
	}

	isInteractive := isInteractive(cfg.DisableProgress)
	cacheBasePath, err := pkgdata.GetCachePath()
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
	err = renderOutput(allPkgs, cfg)
	if err != nil {
		return err
	}

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

func renderOutput(pkgs []*pkgdata.PkgInfo, cfg *config.Config) error {
	switch cfg.OutputFormat {
	case consts.OutputTable:
		out.RenderTable(pkgs, cfg.Fields, cfg.ShowFullTimestamp, cfg.HasNoHeaders)
	case consts.OutputKeyValue:
		out.RenderKeyValue(pkgs, cfg.Fields)
	case consts.OutputJSON:
		out.RenderJSON(pkgs, cfg.Fields)
	default:
		return errors.New("invalid output format")
	}

	return nil
}

func isInteractive(disableProgress bool) bool {
	return term.IsTerminal(int(os.Stdout.Fd())) && !disableProgress
}

func forkCacheWorker(originName string) error {
	cmd := exec.Command(os.Args[0], "--internal-cache-worker="+originName)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start background cache process: %w", err)
	}

	return nil
}

// TODO: add debouncing
func rebuildCache(originName string) error {
	cacheBasePath, err := pkgdata.GetCachePath()
	if err != nil {
		return fmt.Errorf("failed to set up cache dir: %v", err)
	}

	drivers := origins.AvailableDrivers()
	if len(drivers) == 0 {
		return fmt.Errorf("no supported package origins detected")
	}

	var wg sync.WaitGroup
	cfg := &config.Config{
		NoCache:    false,
		RegenCache: false,
	}

	for _, d := range drivers {
		if originName != "all" && originName != d.Name() {
			continue
		}

		wg.Add(1)
		go func(targetDriver driver.Driver) {
			defer wg.Done()
			pipeline := phase.NewPipeline(targetDriver, cfg, false, cacheBasePath)

			if pkgdata.IsLockFileExists(pipeline.CacheRoot) {
				return
			}

			if err := pkgdata.CreateLockFile(pipeline.CacheRoot); err != nil {
				return
			}

			defer pkgdata.RemoveLockFile(pipeline.CacheRoot)

			_, err := pipeline.RunCacheOnly()
			if err != nil {
				return
			}
		}(d)
	}

	wg.Wait()
	return nil
}
