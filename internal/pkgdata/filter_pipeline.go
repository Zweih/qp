package pkgdata

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pipeline/meta"
	"sync"
)

type Filter func(*PkgInfo) bool

type FilterCondition struct {
	Filter    Filter
	PhaseName string
	FieldType consts.FieldType
}

func FilterPackages(
	pkgPtrs []*PkgInfo,
	filterConditions []*FilterCondition,
	reportProgress meta.ProgressReporter,
) []*PkgInfo {
	if len(filterConditions) < 1 {
		return pkgPtrs
	}

	inputChan := populateInitialInputChannel(pkgPtrs)
	outputChan := applyFilterPipeline(inputChan, filterConditions, reportProgress)
	return collectFilteredResults(outputChan)
}

func collectFilteredResults(outputChan <-chan *PkgInfo) []*PkgInfo {
	var filteredPkgPtrs []*PkgInfo

	for pkg := range outputChan {
		filteredPkgPtrs = append(filteredPkgPtrs, pkg)
	}

	return filteredPkgPtrs
}

func applyFilterPipeline(
	inputChan <-chan *PkgInfo,
	filterConditions []*FilterCondition,
	reportProgress meta.ProgressReporter,
) <-chan *PkgInfo {
	outputChan := inputChan
	totalPhases := len(filterConditions)
	completedPhases := 0
	chunkSize := 20

	chunkPool := sync.Pool{
		New: func() any {
			slice := make([]*PkgInfo, 0, chunkSize)
			return &slice
		},
	}

	for filterIndex, f := range filterConditions {
		nextOutputChan := make(chan *PkgInfo, chunkSize)

		go func(inChan <-chan *PkgInfo, outChan chan<- *PkgInfo, filter Filter, phaseName string) {
			defer close(outChan)

			chunkPtr := chunkPool.Get().(*[]*PkgInfo)
			chunk := *chunkPtr
			chunk = chunk[:0]

			for pkg := range inChan {
				chunk = append(chunk, pkg)

				if len(chunk) >= chunkSize {
					processChunk(chunk, outChan, filter)
					chunk = chunk[:0]
				}
			}

			if len(chunk) > 0 {
				processChunk(chunk, outChan, filter)
			}

			*chunkPtr = chunk[:0]
			chunkPool.Put(chunkPtr)

			if reportProgress != nil {
				completedPhases++
				reportProgress(
					completedPhases,
					totalPhases,
					fmt.Sprintf("%s - Step %d/%d completed", phaseName, filterIndex+1, totalPhases),
				)
			}
		}(outputChan, nextOutputChan, f.Filter, f.PhaseName)

		outputChan = nextOutputChan
	}

	return outputChan
}

func processChunk(pkgPtrs []*PkgInfo, outChan chan<- *PkgInfo, filter Filter) {
	for i := range pkgPtrs {
		if filter(pkgPtrs[i]) {
			outChan <- pkgPtrs[i]
		}
	}
}

func populateInitialInputChannel(pkgPtrs []*PkgInfo) <-chan *PkgInfo {
	inputChan := make(chan *PkgInfo, len(pkgPtrs))

	go func() {
		for _, pkg := range pkgPtrs {
			inputChan <- pkg
		}

		close(inputChan)
	}()

	return inputChan
}
