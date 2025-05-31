package phase

import (
	"fmt"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pipeline/meta"
	"qp/internal/pkgdata"
	"qp/internal/quipple/compiler"
	"qp/internal/storage"
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

	cacheMtime, err := storage.LoadCacheModTime(p.CacheRoot)
	if err != nil {
		return nil, nil
	}

	isStale, err := p.Origin.IsCacheStale(cacheMtime)
	if isStale {
		p.ModTime = now
		return nil, err
	}

	pkgs, err := p.Origin.LoadCache(p.CacheRoot)
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

	pkgs, err := p.Origin.Load(p.CacheRoot)
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
	if cfg.NoCache || p.UsedCache {
		return pkgs, nil
	}

	err := p.Origin.SaveCache(p.CacheRoot, pkgs)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: failed to save cache: %v", err))
		return pkgs, nil
	}

	err = storage.SaveCacheModTime(p.CacheRoot, p.ModTime)
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

	return pkgs, nil
}
