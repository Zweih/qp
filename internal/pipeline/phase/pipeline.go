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
	Origin    driver.Driver
	Config    *config.Config
	Pkgs      []*pkgdata.PkgInfo
	UsedCache bool
	CachePath string
	ModTime   int64
}

func NewPipeline(
	origin driver.Driver,
	cfg *config.Config,
	baseCacheDir string,
) *Pipeline {
	cachePath := filepath.Join(baseCacheDir, origin.Name())

	return &Pipeline{
		Origin:    origin,
		Config:    cfg,
		CachePath: cachePath,
	}
}

// TODO: we can merge the logic of these two and just swap in arrays
func (p *Pipeline) Run() ([]*pkgdata.PkgInfo, error) {
	var wg sync.WaitGroup
	phases := []PipelinePhase{
		NewPhase("Load cache", p.loadCacheStep, &wg),
		NewPhase("Fetch packages", p.fetchStep, &wg),
		NewPhase("Resolve dependencies", p.resolveStep, &wg),
		NewPhase("Save cache", p.saveCacheStep, &wg),
	}

	pkgs := []*pkgdata.PkgInfo{}
	var err error

	for _, ph := range phases {
		pkgs, err = ph.Run(p.Config, pkgs)
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

	pkgs := []*pkgdata.PkgInfo{}
	var err error

	for _, ph := range phases {
		pkgs, err = ph.Run(p.Config, pkgs)
		if err != nil {
			return nil, fmt.Errorf("[%s] %w", ph.name, err)
		}
	}

	return pkgs, nil
}
