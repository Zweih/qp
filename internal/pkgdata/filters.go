package pkgdata

import (
	"sync"
	"time"
)

type Filter func([]PackageInfo) []PackageInfo

type FilterCondition struct {
	Condition bool
	Filter    Filter
	PhaseName string
}

type ProgressReporter func(current int, total int, phase string)

func FilterExplicit(pkgs []PackageInfo) []PackageInfo {
	var explicitPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Reason == "explicit" {
			explicitPackages = append(explicitPackages, pkg)
		}
	}

	return explicitPackages
}

func FilterDependencies(pkgs []PackageInfo) []PackageInfo {
	var dependencyPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Reason == "dependency" {
			dependencyPackages = append(dependencyPackages, pkg)
		}
	}

	return dependencyPackages
}

// filters packages installed on specific date
func FilterByDate(pkgs []PackageInfo, date time.Time) []PackageInfo {
	var filteredPackages []PackageInfo

	for _, pkg := range pkgs {
		if pkg.Timestamp.Year() == date.Year() && pkg.Timestamp.YearDay() == date.YearDay() {
			filteredPackages = append(filteredPackages, pkg)
		}
	}

	return filteredPackages
}

func FilterBySize(pkgs []PackageInfo, operator string, sizeInBytes int64) []PackageInfo {
	var filteredPackages []PackageInfo

	for _, pkg := range pkgs {
		switch operator {
		case ">":
			if pkg.Size > sizeInBytes {
				filteredPackages = append(filteredPackages, pkg)
			}
		case "<":
			if pkg.Size < sizeInBytes {
				filteredPackages = append(filteredPackages, pkg)
			}
		}
	}

	return filteredPackages
}

func applyConcurrentFilter(packages []PackageInfo, filterFunc Filter) []PackageInfo {
	const chunkSize = 100

	var mu sync.Mutex
	var wg sync.WaitGroup
	var filteredPackages []PackageInfo

	for i := 0; i < len(packages); i += chunkSize {
		endIdx := i + chunkSize

		if endIdx > len(packages) {
			endIdx = len(packages)
		}

		chunk := packages[i:endIdx]

		wg.Add(1)

		go func(chunk []PackageInfo) {
			defer wg.Done()

			filteredChunk := filterFunc(chunk)

			mu.Lock()
			filteredPackages = append(filteredPackages, filteredChunk...)
			mu.Unlock()
		}(chunk)
	}

	wg.Wait()

	return filteredPackages
}

func ApplyFilters(
	pkgs []PackageInfo,
	filters []FilterCondition,
	reportProgress ProgressReporter,
) []PackageInfo {
	totalFilters := len(filters)
	currentFilter := 0

	for _, f := range filters {
		if f.Condition {
			if reportProgress != nil {
				reportProgress(currentFilter, totalFilters, f.PhaseName)
			}

			pkgs = applyConcurrentFilter(pkgs, f.Filter)

			if reportProgress != nil {
				currentFilter++
				reportProgress(currentFilter, totalFilters, f.PhaseName)
			}
		}
	}

	if reportProgress != nil {
		reportProgress(totalFilters, totalFilters, "All filters completed")
	}

	return pkgs
}
