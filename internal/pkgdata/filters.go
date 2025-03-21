package pkgdata

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Filter func(*PkgInfo) bool

type FilterCondition struct {
	Filter    Filter
	PhaseName string
}

func FilterByRelation(pkgNames []string, targetNames []string) bool {
	for _, targetName := range targetNames {
		for _, packageName := range pkgNames {
			matches := packageNameRegex.FindStringSubmatch(packageName)
			if len(matches) >= 2 && matches[1] == targetName {
				return true
			}
		}
	}

	return false
}

func FilterByReason(installReason string, targetReason string) bool {
	return installReason == targetReason
}

func FilterExplicit(pkg *PkgInfo) bool {
	return pkg.Reason == "explicit"
}

func FilterDependencies(pkg *PkgInfo) bool {
	return pkg.Reason == "dependency"
}

// filters for packages installed on specific date
func FilterByDate(pkg *PkgInfo, date int64) bool {
	pkgDate := time.Unix(pkg.Timestamp, 0)
	targetDate := time.Unix(date, 0)
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
	for _, targetString := range targetStrings {
		if strings.Contains(pkgString, targetString) {
			return true
		}
	}

	return false
}

func FilterPackages(
	pkgs []PkgInfo,
	filterConditions []FilterCondition,
	reportProgress ProgressReporter,
) []PkgInfo {
	if len(filterConditions) < 1 {
		return pkgs
	}

	pkgPtrs := convertToPointers(pkgs)

	inputChan := populateInitialInputChannel(pkgPtrs)
	outputChan := applyFilterPipeline(inputChan, filterConditions, reportProgress)
	return collectFilteredResults(outputChan)
}

func collectFilteredResults(outputChan <-chan *PkgInfo) []PkgInfo {
	var filteredPkgs []PkgInfo

	for pkg := range outputChan {
		filteredPkgs = append(filteredPkgs, *pkg)
	}

	return filteredPkgs
}

func applyFilterPipeline(
	inputChan <-chan *PkgInfo,
	filterConditions []FilterCondition,
	reportProgress ProgressReporter,
) <-chan *PkgInfo {
	outputChan := inputChan
	totalPhases := len(filterConditions)
	completedPhases := 0
	chunkSize := 100

	for filterIndex, f := range filterConditions {
		nextOutputChan := make(chan *PkgInfo, cap(inputChan))

		go func(inChan <-chan *PkgInfo, outChan chan<- *PkgInfo, filter Filter, phaseName string) {
			defer close(outChan)

			var chunk []*PkgInfo
			for pkg := range inChan {
				chunk = append(chunk, pkg)

				if len(chunk) >= chunkSize {
					processChunk(chunk, outChan, filter)
					chunk = nil
				}
			}

			if len(chunk) > 0 {
				processChunk(chunk, outChan, filter)
			}

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
