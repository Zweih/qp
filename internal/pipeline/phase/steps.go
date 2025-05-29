package phase

import (
	"fmt"
	"path/filepath"
	"qp/internal/compiler"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pipeline/filtering"
	"qp/internal/pipeline/meta"
	"qp/internal/pkgdata"
	"time"
)

func (p *Pipeline) loadCacheStep(
	cfg *config.Config,
	_ []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	now := time.Now().Unix()
	if cfg.RegenCache {
		p.ModTime = now
	}

	if cfg.RegenCache || cfg.NoCache {
		return nil, nil
	}

	cacheRoot := filepath.Join(p.CachePath, p.Origin.Name())
	cacheMtime, err := pkgdata.LoadCacheModTime(cacheRoot)
	if err != nil {
		return nil, nil
	}

	isStale, err := p.Origin.IsCacheStale(cacheMtime)
	if isStale {
		p.ModTime = now
		return nil, err
	}

	pkgs, err := p.Origin.LoadCache(cacheRoot)
	if err != nil {
		p.ModTime = now
		return nil, nil
	}

	p.UsedCache = true
	p.ModTime = cacheMtime

	return pkgs, nil
}

func (p *Pipeline) fetchStep(
	_ *config.Config,
	pkgs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if p.UsedCache {
		return pkgs, nil
	}

	pkgs, err := p.Origin.Load()
	if err != nil {
		err = fmt.Errorf("failed to fetch packages: %v", err)
		return nil, err
	}

	return pkgs, nil
}

func (p *Pipeline) resolveStep(
	_ *config.Config,
	pkgs []*pkgdata.PkgInfo,
	reportProgress meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if p.UsedCache {
		return pkgs, nil
	}

	pkgs, err := p.Origin.ResolveDeps(pkgs)
	if err != nil {
		return nil, fmt.Errorf("dependency resolution failed: %w", err)
	}

	return pkgs, nil
}

func (p *Pipeline) saveCacheStep(
	cfg *config.Config,
	pkgs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	cacheRoot := filepath.Join(p.CachePath, p.Origin.Name())
	pkgdata.UpdateInstallHistory(cacheRoot, pkgs)

	if cfg.NoCache || p.UsedCache {
		return pkgs, nil
	}

	err := p.Origin.SaveCache(cacheRoot, pkgs)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: failed to save cache: %v", err))
		return pkgs, nil
	}

	err = pkgdata.SaveCacheModTime(cacheRoot, p.ModTime)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: failed to save cache modtime: %v", err))
		return pkgs, nil
	}

	return pkgs, nil
}

func (p *Pipeline) filterStep(
	cfg *config.Config,
	pkgs []*pkgdata.PkgInfo,
	reportProgress meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if cfg.QueryExpr != nil {
		return compiler.RunDAG(cfg.QueryExpr, pkgs)
	}

	if len(cfg.FieldQueries) == 0 {
		return pkgs, nil
	}

	filterConditions, err := filtering.QueriesToConditions(cfg.FieldQueries)
	if err != nil {
		return nil, fmt.Errorf("filter query error: %w", err)
	}

	return pkgdata.FilterPackages(pkgs, filterConditions, reportProgress), nil
}
