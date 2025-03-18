package pkgdata

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Filter func(PackageInfo) bool

type FilterCondition struct {
	Filter    Filter
	PhaseName string
}

func FilterByRelation(packageNames []string, targetNames []string) bool {
	for _, targetName := range targetNames {
		for _, packageName := range packageNames {
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

func FilterExplicit(pkg PackageInfo) bool {
	return pkg.Reason == "explicit"
}

func FilterDependencies(pkg PackageInfo) bool {
	return pkg.Reason == "dependency"
}

// // filters for packages installed on specific date
func FilterByDate(pkg PackageInfo, date int64) bool {
	pkgDate := time.Unix(pkg.Timestamp, 0)
	targetDate := time.Unix(date, 0)
	return pkgDate.Year() == targetDate.Year() && pkgDate.YearDay() == targetDate.YearDay()
}

// inclusive
func FilterByDateRange(pkg PackageInfo, start int64, end int64) bool {
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
func FilterBySize(pkg PackageInfo, targetSize int64) bool {
	return roundSizeInBytes(pkg.Size) == roundSizeInBytes(targetSize)
}

func FilterBySizeRange(pkg PackageInfo, startSize int64, endSize int64) bool {
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
	pkgs []PackageInfo,
	filters []FilterCondition,
	reportProgress ProgressReporter,
) []PackageInfo {
	if len(filters) < 1 {
		return pkgs
	}

	inputChan := populateInitialInputChannel(pkgs)
	outputChan := applyFilterPipeline(inputChan, filters, reportProgress)
	return collectFilteredResults(outputChan)
}

func collectFilteredResults(outputChan <-chan PackageInfo) []PackageInfo {
	var filteredPackages []PackageInfo

	for pkg := range outputChan {
		filteredPackages = append(filteredPackages, pkg)
	}

	return filteredPackages
}

func applyFilterPipeline(
	inputChan <-chan PackageInfo,
	filters []FilterCondition,
	reportProgress ProgressReporter,
) <-chan PackageInfo {
	outputChan := inputChan
	totalPhases := len(filters)
	completedPhases := 0

	for filterIndex, f := range filters {
		nextOutputChan := make(chan PackageInfo, cap(inputChan))

		go func(
			inChan <-chan PackageInfo,
			outChan chan<- PackageInfo,
			filter Filter,
			phaseName string,
		) {
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

func populateInitialInputChannel(pkgs []PackageInfo) <-chan PackageInfo {
	inputChan := make(chan PackageInfo, len(pkgs))

	go func() {
		for _, pkg := range pkgs {
			inputChan <- pkg
		}

		close(inputChan)
	}()

	return inputChan
}
