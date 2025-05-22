package phase

import (
	"fmt"
	"path/filepath"
	"qp/api/driver"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pkgdata"
	"sync"
	"time"
)

type Pipeline struct {
	Origin        driver.Driver
	Config        *config.Config
	Pkgs          []*pkgdata.PkgInfo
	IsInteractive bool
	UsedCache     bool
	CachePath     string
	ModTime       int64
}

func NewPipeline(
	origin driver.Driver,
	cfg *config.Config,
	isInteractive bool,
	baseCachePath string,
) *Pipeline {
	modTime, err := origin.SourceModified()
	if err != nil {
		out.WriteLine(fmt.Sprintf("WARNING: failed to get mod time for origin %s: %v", origin.Name(), err))
		modTime = time.Now().Unix()
	}

	return &Pipeline{
		Origin:        origin,
		Config:        cfg,
		IsInteractive: isInteractive,
		CachePath:     filepath.Join(baseCachePath, origin.Name()+".cache"),
		ModTime:       modTime,
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
