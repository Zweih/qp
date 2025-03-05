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

func FilterRequiredBy(pkg PackageInfo, target string) bool {
	for _, requiredByPkgName := range pkg.RequiredBy {
		matches := packageNameRegex.FindStringSubmatch(requiredByPkgName)
		if len(matches) >= 2 {
			if matches[1] == target { // TODO: condense
				return true
			}
		}
	}

	return false
}

func FilterExplicit(pkg PackageInfo) bool {
	return pkg.Reason == "explicit"
}

func FilterDependencies(pkg PackageInfo) bool {
	return pkg.Reason == "dependency"
}

// filters for packages installed on specific date
func FilterByDate(pkg PackageInfo, date time.Time) bool {
	return pkg.Timestamp.Year() == date.Year() && pkg.Timestamp.YearDay() == date.YearDay()
}

// inclusive
func FilterByDateRange(pkg PackageInfo, startDate time.Time, endDate time.Time) bool {
	return !(pkg.Timestamp.Before(startDate) || pkg.Timestamp.After(endDate))
}

func roundSizeInBytes(num int64) int64 {
	if num < 1000 {
		return num
	}

	numDigits := int(math.Log10(float64(num))) + 1
	scaleFactor := int64(math.Pow10(numDigits - 3))

	return num / scaleFactor
}

func FilterBySize(pkg PackageInfo, size int64) bool {
	return roundSizeInBytes(pkg.Size) == roundSizeInBytes(size)
}

func FilterBySizeRange(pkg PackageInfo, startSize int64, endSize int64) bool {
	return pkg.Size >= startSize && pkg.Size <= endSize
}

func FilterByName(pkg PackageInfo, searchTerm string) bool {
	return strings.Contains(pkg.Name, searchTerm)
}

func ApplyFilters(pkgs []PackageInfo, filters []FilterCondition, reportProgress ProgressReporter) []PackageInfo {
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
