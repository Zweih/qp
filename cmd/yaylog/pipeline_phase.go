package main

import (
	"fmt"
	"strings"
	"yaylog/internal/config"
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
	go phase.displayProgress(progressChan)

	return progressChan
}

func (phase PipelinePhase) stopProgress(progressChan chan ProgressMessage) {
	if progressChan != nil {
		close(progressChan)
	}
}

func (phase PipelinePhase) displayProgress(progressChan chan ProgressMessage) {
	for msg := range progressChan {
		fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
		progressMessageText := fmt.Sprintf("[%s] %d%% - %s", msg.Phase, msg.Progress, msg.Description)
		fmt.Printf("%-80s", progressMessageText)

		if msg.Progress == 100 {
			// extra clear when done
			fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
		}
	}
}
