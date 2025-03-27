package main

import (
	"fmt"
	"sync"
	"yaylog/internal/config"
	out "yaylog/internal/display"
	"yaylog/internal/pipeline/meta"
	"yaylog/internal/pkgdata"
)

type (
	ProgressReporter = meta.ProgressReporter
	ProgressMessage  = meta.ProgressMessage
	PkgInfo          = pkgdata.PkgInfo
)

type Operation func(
	cfg config.Config,
	packages []*PkgInfo,
	progressReporter meta.ProgressReporter,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error)

type PipelinePhase struct {
	Name      string
	Operation Operation
	wg        *sync.WaitGroup
}

func (phase PipelinePhase) Run(
	cfg config.Config,
	packages []*PkgInfo,
	pipelineCtx *meta.PipelineContext,
) ([]*PkgInfo, error) {
	progressChan := phase.startProgress(pipelineCtx.IsInteractive)
	outputPackages, err := phase.Operation(
		cfg,
		packages,
		phase.reportProgress(progressChan),
		pipelineCtx,
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
			Description: fmt.Sprintf(("%s is in progress..."), phase.Name),
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
