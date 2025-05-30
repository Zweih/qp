package phase

import (
	"fmt"
	"path/filepath"
	"qp/api/driver"
	"qp/internal/config"
	"qp/internal/pkgdata"
	"sync"
)

type Pipeline struct {
	Origin        driver.Driver
	Config        *config.Config
	Pkgs          []*pkgdata.PkgInfo
	IsInteractive bool
	UsedCache     bool
	CacheRoot     string
	ModTime       int64
}

func NewPipeline(
	origin driver.Driver,
	cfg *config.Config,
	isInteractive bool,
	baseCachePath string,
) *Pipeline {
	cacheRoot := filepath.Join(baseCachePath, origin.Name())

	return &Pipeline{
		Origin:        origin,
		Config:        cfg,
		IsInteractive: isInteractive,
		CacheRoot:     cacheRoot,
	}
}

func (p *Pipeline) Run() ([]*pkgdata.PkgInfo, error) {
	var wg sync.WaitGroup
	phases := []PipelinePhase{
		NewPhase("Load cache", p.loadCacheStep, &wg),
		NewPhase("Fetch packages", p.fetchStep, &wg),
		NewPhase("Resolve dependencies", p.resolveStep, &wg),
		NewPhase("Save cache", p.saveCacheStep, &wg),
		NewPhase("Filter", p.filterStep, &wg),
	}

	pkgs := []*PkgInfo{}
	var err error

	for _, ph := range phases {
		pkgs, err = ph.Run(p.Config, pkgs, p.IsInteractive)
		if err != nil {
			return nil, fmt.Errorf("[%s] %w", ph.name, err)
		}
	}

	return pkgs, nil
}

func (p *Pipeline) RunCacheOnly() ([]*pkgdata.PkgInfo, error) {
	var wg sync.WaitGroup
	phases := []PipelinePhase{
		NewPhase("Fetch packages", p.fetchStep, &wg),
		NewPhase("Resolve dependencies", p.resolveStep, &wg),
		NewPhase("Save cache", p.saveCacheStep, &wg),
	}

	pkgs := []*PkgInfo{}
	var err error

	for _, ph := range phases {
		pkgs, err = ph.Run(p.Config, pkgs, p.IsInteractive)
		if err != nil {
			return nil, fmt.Errorf("[%s] %w", ph.name, err)
		}
	}

	return pkgs, nil
}
