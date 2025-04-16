package pkgdata

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/pipeline/meta"
	"slices"
	"strings"
	"sync"
	"time"
)

type Filter func(*PkgInfo) bool

type FilterCondition struct {
	Filter    Filter
	PhaseName string
	FieldType consts.FieldType
}

const fuzzySizeTolerancePercent = 0.3

func FilterByReason(installReason string, targetReason string) bool {
	return installReason == targetReason
}

func GetRelationsByDepth(relations []Relation, targetDepth int32) []Relation {
	filteredRelations := []Relation{}

	for _, relation := range relations {
		if relation.Depth == targetDepth {
			filteredRelations = append(filteredRelations, relation)
		}
	}

	return filteredRelations
}

func FuzzyDate(pkg *PkgInfo, date int64) bool {
	pkgDate := time.Unix(pkg.InstallTimestamp, 0)
	targetDate := time.Unix(date, 0) // TODO: we can pull this out to the top level
	return pkgDate.Year() == targetDate.Year() && pkgDate.YearDay() == targetDate.YearDay()
}

func StrictDateRange(pkg *PkgInfo, start int64, end int64) bool {
	return !(pkg.InstallTimestamp < start || pkg.InstallTimestamp > end)
}

func StrictDate(pkg *PkgInfo, targetDate int64) bool {
	return pkg.InstallTimestamp == targetDate
}

func FuzzyDateRange(pkg *PkgInfo, start int64, end int64) bool {
	pkgDate := time.Unix(pkg.InstallTimestamp, 0).Truncate(24 * time.Hour)
	startDate := time.Unix(start, 0).Truncate(24 * time.Hour)
	endDate := time.Unix(end, 0).Truncate(24 * time.Hour)

	return (pkgDate.Equal(startDate) || pkgDate.After(startDate)) &&
		(pkgDate.Equal(endDate) || pkgDate.Before(endDate))
}

func FuzzySizeTolerance(targetSize int64) int64 {
	return int64(float64(targetSize) * fuzzySizeTolerancePercent / 100.0)
}

func FuzzySize(pkg *PkgInfo, targetSize int64) bool {
	tolerance := FuzzySizeTolerance(targetSize)
	result := pkg.Size - targetSize
	return max(result, -result) <= tolerance
}

func FuzzySizeRange(pkg *PkgInfo, start int64, end int64) bool {
	toleranceStart := FuzzySizeTolerance(start)
	toleranceEnd := FuzzySizeTolerance(end)

	return pkg.Size >= (start-toleranceStart) && pkg.Size <= (end+toleranceEnd)
}

func StrictSize(pkg *PkgInfo, targetSize int64) bool {
	return pkg.Size == targetSize
}

func StrictSizeRange(pkg *PkgInfo, startSize int64, endSize int64) bool {
	return !(pkg.Size < startSize || pkg.Size > endSize)
}

func FilterSliceByStrings(pkgStrings []string, targetStrings []string) bool {
	for _, pkgString := range pkgStrings {
		if FuzzyStrings(pkgString, targetStrings) {
			return true
		}
	}

	return false
}

func FuzzyStrings(pkgString string, targetStrings []string) bool {
	pkgString = strings.ToLower(pkgString)

	for _, targetString := range targetStrings {
		if strings.Contains(pkgString, targetString) {
			return true
		}
	}

	return false
}

func StrictStrings(pkgString string, targetStrings []string) bool {
	pkgString = strings.ToLower(pkgString)

	return slices.Contains(targetStrings, pkgString)
}

func RelationExists(relations []Relation) bool {
	return len(relations) > 0
}

func StringExists(pkgString string) bool {
	return pkgString != ""
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
