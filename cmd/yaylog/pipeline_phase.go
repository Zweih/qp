package main

import (
	"fmt"
	"sync"
	"yaylog/internal/config"
	"yaylog/internal/display"
	"yaylog/internal/pkgdata"
)

type (
	ProgressReporter = pkgdata.ProgressReporter
	ProgressMessage  = pkgdata.ProgressMessage
	PackageInfo      = pkgdata.PackageInfo
)

type PipelinePhase struct {
	Name          string
	Operation     func(cfg config.Config, packages []PackageInfo, progressReporter ProgressReporter) []PackageInfo
	IsInteractive bool
	wg            *sync.WaitGroup
}

func (phase PipelinePhase) Run(cfg config.Config, packages []PackageInfo) []PackageInfo {
	progressChan := phase.startProgress()
	outputPackages := phase.Operation(cfg, packages, phase.reportProgress(progressChan))
	phase.stopProgress(progressChan)

	return outputPackages
}

func (phase PipelinePhase) reportProgress(progressChan chan ProgressMessage) ProgressReporter {
	if progressChan == nil {
		return ProgressReporter(func(current int, total int, phaseName string) {})
	}

	return ProgressReporter(func(current int, total int, phaseName string) {
		progressChan <- ProgressMessage{
			Phase:       phaseName,
			Progress:    (current * 100) / total,
			Description: fmt.Sprintf(("%s is in progress..."), phase.Name),
		}
	})
}

func (phase PipelinePhase) startProgress() chan ProgressMessage {
	if !phase.IsInteractive {
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
		display.Manager.ClearProgress()
	}
}

func (phase PipelinePhase) displayProgress(progressChan chan ProgressMessage) {
	for msg := range progressChan {
		display.Manager.PrintProgress(msg.Phase, msg.Progress, msg.Description)
	}

	display.Manager.ClearProgress()
}
