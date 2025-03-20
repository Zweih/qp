package pkgdata

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Filter func(PkgInfo) bool

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

func FilterExplicit(pkg PkgInfo) bool {
	return pkg.Reason == "explicit"
}

func FilterDependencies(pkg PkgInfo) bool {
	return pkg.Reason == "dependency"
}

// filters for packages installed on specific date
func FilterByDate(pkg PkgInfo, date int64) bool {
	pkgDate := time.Unix(pkg.Timestamp, 0)
	targetDate := time.Unix(date, 0)
	return pkgDate.Year() == targetDate.Year() && pkgDate.YearDay() == targetDate.YearDay()
}

// inclusive
func FilterByDateRange(pkg PkgInfo, start int64, end int64) bool {
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
func FilterBySize(pkg PkgInfo, targetSize int64) bool {
	return roundSizeInBytes(pkg.Size) == roundSizeInBytes(targetSize)
}

func FilterBySizeRange(pkg PkgInfo, startSize int64, endSize int64) bool {
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
	filters []FilterCondition,
	reportProgress ProgressReporter,
) []PkgInfo {
	if len(filters) < 1 {
		return pkgs
	}

	inputChan := populateInitialInputChannel(pkgs)
	outputChan := applyFilterPipeline(inputChan, filters, reportProgress)
	return collectFilteredResults(outputChan)
}

func collectFilteredResults(outputChan <-chan PkgInfo) []PkgInfo {
	var filteredPackages []PkgInfo

	for pkg := range outputChan {
		filteredPackages = append(filteredPackages, pkg)
	}

	return filteredPackages
}

func applyFilterPipeline(
	inputChan <-chan PkgInfo,
	filters []FilterCondition,
	reportProgress ProgressReporter,
) <-chan PkgInfo {
	outputChan := inputChan
	totalPhases := len(filters)
	completedPhases := 0

	for filterIndex, f := range filters {
		nextOutputChan := make(chan PkgInfo, cap(inputChan))

		go func(inChan <-chan PkgInfo, outChan chan<- PkgInfo, filter Filter, phaseName string) {
			for pkg := range inChan {
				if filter(pkg) {
					outChan <- pkg
				}
			}

			if reportProgress != nil {
				completedPhases++
				reportProgress(
					completedPhases,
					totalPhases,
					fmt.Sprintf("%s - Step %d/%d completed", phaseName, filterIndex+1, totalPhases),
				)
			}

			close(outChan)
		}(outputChan, nextOutputChan, f.Filter, f.PhaseName)

		outputChan = nextOutputChan
	}

	return outputChan
}

func populateInitialInputChannel(pkgs []PkgInfo) <-chan PkgInfo {
	inputChan := make(chan PkgInfo, len(pkgs))

	go func() {
		for _, pkg := range pkgs {
			inputChan <- pkg
		}

		close(inputChan)
	}()

	return inputChan
}
