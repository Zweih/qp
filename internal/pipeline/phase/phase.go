package phase

import (
	"qp/internal/config"
	"qp/internal/pkgdata"
	"sync"
)

type Step func(
	cfg *config.Config,
	pkgs []*pkgdata.PkgInfo,
) ([]*pkgdata.PkgInfo, error)

type PipelinePhase struct {
	name string
	step Step
	wg   *sync.WaitGroup
}

func NewPhase(name string, step Step, wg *sync.WaitGroup) PipelinePhase {
	return PipelinePhase{name, step, wg}
}

func (phase PipelinePhase) Run(
	cfg *config.Config,
	pkgs []*pkgdata.PkgInfo,
) ([]*pkgdata.PkgInfo, error) {
	outputPackages, err := phase.step(cfg, pkgs)
	if err != nil {
		return nil, err
	}

	return outputPackages, nil
}
