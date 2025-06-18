package phase

import (
	"fmt"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pkgdata"
	"qp/internal/storage"
	"time"
)

func (p *Pipeline) loadCacheStep(
	cfg *config.Config,
	_ []*pkgdata.PkgInfo,
) ([]*pkgdata.PkgInfo, error) {
	now := time.Now().Unix()
	if cfg.RegenCache {
		p.ModTime = now
	}

	if cfg.RegenCache || cfg.NoCache {
		return nil, nil
	}

	cacheMtime, err := storage.LoadCacheModTime(p.CachePath)
	if err != nil {
		return nil, nil
	}

	isStale, err := p.Origin.IsCacheStale(cacheMtime)
	if isStale {
		p.ModTime = now
		return nil, err
	}

	pkgs, err := p.Origin.LoadCache(p.CachePath)
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
) ([]*pkgdata.PkgInfo, error) {
	if p.UsedCache {
		return pkgs, nil
	}

	p.ModTime = time.Now().Unix()

	pkgs, err := p.Origin.Load(p.CachePath)
	if err != nil {
		err = fmt.Errorf("failed to fetch packages: %v", err)
		return nil, err
	}

	return pkgs, nil
}

func (p *Pipeline) resolveStep(
	_ *config.Config,
	pkgs []*pkgdata.PkgInfo,
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
) ([]*pkgdata.PkgInfo, error) {
	if cfg.NoCache || p.UsedCache {
		return pkgs, nil
	}

	err := p.Origin.SaveCache(p.CachePath, pkgs)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: failed to save cache: %v", err))
		return pkgs, nil
	}

	err = storage.SaveCacheModTime(p.CachePath, p.ModTime)
	if err != nil {
		out.WriteLine(fmt.Sprintf("Warning: failed to save cache modtime: %v", err))
		return pkgs, nil
	}

	return pkgs, nil
}
