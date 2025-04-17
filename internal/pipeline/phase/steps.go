package phase

import (
	"fmt"
	"path/filepath"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pipeline/filtering"
	"qp/internal/pipeline/meta"
	"qp/internal/pkgdata"
)

func (p *Pipeline) loadCacheStep(
	cfg *config.Config,
	_ []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if cfg.RegenCache || cfg.NoCache {
		return nil, nil
	}

	cachePath := filepath.Join(p.CachePath)
	pkgs, err := p.Origin.LoadCache(cachePath, p.ModTime)
	if err == nil {
		p.UsedCache = true
	}

	return pkgs, nil
}

func (p *Pipeline) fetchStep(
	_ *config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if p.UsedCache {
		return pkgPtrs, nil
	}

	pkgs, err := p.Origin.Load()
	if err != nil {
		out.WriteLine(fmt.Sprintf(
			"Warning: failed to fetch packages for origin %s: %v",
			p.Origin.Name(), err,
		))
		return nil, err
	}

	return pkgs, nil
}

func (p *Pipeline) resolveStep(
	_ *config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	reportProgress meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if p.UsedCache {
		return pkgPtrs, nil
	}

	pkgs, err := p.Origin.ResolveDeps(pkgPtrs)
	if err != nil {
		return nil, fmt.Errorf("dependency resolution failed for origin %s: %w", p.Origin.Name(), err)
	}

	return pkgs, nil
}

func (p *Pipeline) saveCacheStep(
	cfg *config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	_ meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if cfg.NoCache || p.UsedCache {
		return pkgPtrs, nil
	}

	cachePath := filepath.Join(p.CachePath)
	err := p.Origin.SaveCache(cachePath, pkgPtrs, p.ModTime)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: failed to save cache for origin %s: %v", p.Origin.Name(), err))
	}

	return pkgPtrs, nil
}

func (p *Pipeline) filterStep(
	cfg *config.Config,
	pkgPtrs []*pkgdata.PkgInfo,
	reportProgress meta.ProgressReporter,
) ([]*pkgdata.PkgInfo, error) {
	if len(cfg.FieldQueries) == 0 {
		return pkgPtrs, nil
	}

	filterConditions, err := filtering.QueriesToConditions(cfg.FieldQueries)
	if err != nil {
		return nil, fmt.Errorf("filter query error: %w", err)
	}

	return pkgdata.FilterPackages(pkgPtrs, filterConditions, reportProgress), nil
}
