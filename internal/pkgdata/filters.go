package pkgdata

import (
	"fmt"
	"math"
	"qp/internal/consts"
	"qp/internal/pipeline/meta"
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

func FilterByRelation(relations []Relation, targetNames []string) bool {
	for _, targetName := range targetNames {
		for _, relation := range relations {
			if relation.Name == targetName {
				return true
			}
		}
	}

	return false
}

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

// filters for packages installed on specific date
func FilterByDate(pkg *PkgInfo, date int64) bool {
	pkgDate := time.Unix(pkg.Timestamp, 0)
	targetDate := time.Unix(date, 0) // TODO: we can pull this out to the top level
	return pkgDate.Year() == targetDate.Year() && pkgDate.YearDay() == targetDate.YearDay()
}

// inclusive
func FilterByDateRange(pkg *PkgInfo, start int64, end int64) bool {
	return !(pkg.Timestamp < start || pkg.Timestamp > end)
}

func roundSizeInBytes(num int64) int64 {
	if num < 1000 {
		return num
	}

	numDigits := int(math.Log10(float64(num))) + 1
	scaleFactor := int64(math.Pow10(numDigits - 3))

	return num / scaleFactor
}

// TODO: let's pre-round the inputs outside of these functions
func FilterBySize(pkg *PkgInfo, targetSize int64) bool {
	return roundSizeInBytes(pkg.Size) == roundSizeInBytes(targetSize)
}

func FilterBySizeRange(pkg *PkgInfo, startSize int64, endSize int64) bool {
	roundedSize := roundSizeInBytes(pkg.Size)
	return !(roundedSize < roundSizeInBytes(startSize) || roundedSize > roundSizeInBytes(endSize))
}

func FilterByStrings(pkgString string, targetStrings []string) bool {
	pkgString = strings.ToLower(pkgString)

	for _, targetString := range targetStrings {
		if strings.Contains(pkgString, targetString) {
			return true
		}
	}

	return false
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
