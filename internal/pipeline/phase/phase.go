package phase

import (
	"fmt"
	"qp/internal/config"
	out "qp/internal/display"
	"qp/internal/pipeline/meta"
	"qp/internal/pkgdata"
	"sync"
)

type (
	ProgressReporter = meta.ProgressReporter
	ProgressMessage  = meta.ProgressMessage
	PkgInfo          = pkgdata.PkgInfo
)

type Step func(
	cfg *config.Config,
	packages []*PkgInfo,
	progressReporter meta.ProgressReporter,
) ([]*PkgInfo, error)

type PipelinePhase struct {
	name string
	step Step
	wg   *sync.WaitGroup
}

func NewPhase(name string, step Step, wg *sync.WaitGroup) PipelinePhase {
	return PipelinePhase{
		name,
		step,
		wg,
	}
}

func (phase PipelinePhase) Run(
	cfg *config.Config,
	packages []*PkgInfo,
	isInteractive bool,
) ([]*PkgInfo, error) {
	progressChan := phase.startProgress(isInteractive)
	outputPackages, err := phase.step(
		cfg,
		packages,
		phase.reportProgress(progressChan),
	)
	phase.stopProgress(progressChan)

	if err != nil {
		return nil, err
	}

	return outputPackages, nil
}

func (phase PipelinePhase) reportProgress(progressChan chan ProgressMessage) ProgressReporter {
	if progressChan == nil {
		return ProgressReporter(func(_ int, _ int, _ string) {})
	}

	return ProgressReporter(func(current int, total int, phaseName string) {
		progressChan <- ProgressMessage{
			Phase:       phaseName,
			Progress:    (current * 100) / total,
			Description: fmt.Sprintf(("%s is in progress..."), phase.name),
		}
	})
}

func (phase PipelinePhase) startProgress(isInteractive bool) chan ProgressMessage {
	if !isInteractive {
		return nil
	}

	progressChan := make(chan ProgressMessage)
	phase.wg.Add(1)

	go func() {
		defer phase.wg.Done()
		phase.displayProgress(progressChan)
	}()

	return progressChan
}

func (phase PipelinePhase) stopProgress(progressChan chan ProgressMessage) {
	if progressChan != nil {
		close(progressChan)
		phase.wg.Wait()
		out.ClearProgress()
	}
}

func (phase PipelinePhase) displayProgress(progressChan chan ProgressMessage) {
	for msg := range progressChan {
		out.PrintProgress(msg.Phase, msg.Progress, msg.Description)
	}

	out.ClearProgress()
}
