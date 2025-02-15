package pkgdata

import (
	"strings"
	"time"
)

type Filter func(PackageInfo) bool

type FilterCondition struct {
	Condition bool
	Filter    Filter
	PhaseName string
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

func FilterBySize(pkg PackageInfo, operator string, sizeInBytes int64) bool {
	switch operator {
	case ">":
		return pkg.Size > sizeInBytes
	case "<":
		return pkg.Size < sizeInBytes
	default:
		return false
	}
}

func FilterByName(pkg PackageInfo, searchTerm string) bool {
	return strings.Contains(pkg.Name, searchTerm)
}

func ApplyFilters(pkgs []PackageInfo, filters []FilterCondition, reportProgress ProgressReporter) []PackageInfo {
	if len(filters) < 1 {
		return pkgs
	}

	inputChan := populateInitialInputChannel(pkgs)
	outputChan := applyFilterPipeline(inputChan, filters)
	return collectFilteredResults(outputChan)
}

func collectFilteredResults(outputChan <-chan PackageInfo) []PackageInfo {
	var filteredPackages []PackageInfo

	for pkg := range outputChan {
		filteredPackages = append(filteredPackages, pkg)
	}

	return filteredPackages
}

func applyFilterPipeline(inputChan <-chan PackageInfo, filters []FilterCondition) <-chan PackageInfo {
	outputChan := inputChan

	for _, f := range filters {
		if !f.Condition {
			continue
		}

		nextOutputChan := make(chan PackageInfo, cap(inputChan))

		go func(inChan <-chan PackageInfo, outChan chan<- PackageInfo, filter Filter, phaseName string) {
			for pkg := range inChan {
				if filter(pkg) {
					outChan <- pkg
				}
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
