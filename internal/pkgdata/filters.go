package pkgdata

import (
	"fmt"
	"sync"
	"time"
	"yaylog/internal/config"
)

type FilterCondition struct {
	Condition bool
	Filter    func([]PackageInfo) []PackageInfo
	PhaseName string
}

type ProgressMessage struct {
	Phase       string
	Progress    int
	Description string
}

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

func applyConcurrentFilter(packages []PackageInfo, filterFunc func([]PackageInfo) []PackageInfo) []PackageInfo {
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

func ConcurrentFilters(
	packages []PackageInfo,
	progressChan chan ProgressMessage,
	dateFilter time.Time,
	sizeFilter config.SizeFilter,
	explicitOnly bool,
	dependenciesOnly bool,
) []PackageInfo {
	filters := []FilterCondition{
		{
			Condition: explicitOnly,
			Filter:    FilterExplicit,
			PhaseName: "filtering by explicit packages",
		},
		{
			Condition: dependenciesOnly,
			Filter:    FilterDependencies,
			PhaseName: "filtering by dependencies",
		},
		{
			Condition: !dateFilter.IsZero(),
			Filter: func(pkgs []PackageInfo) []PackageInfo {
				return FilterByDate(pkgs, dateFilter)
			},
			PhaseName: "filtering by date",
		},
		{
			Condition: sizeFilter.IsFilter,
			Filter: func(pkgs []PackageInfo) []PackageInfo {
				return FilterBySize(pkgs, sizeFilter.Operator, sizeFilter.SizeInBytes)
			},
			PhaseName: "filtering by size",
		},
	}

	totalFilters := 0

	for _, f := range filters {
		if f.Condition {
			totalFilters++
		}
	}

	filtersApplied := 0

	for _, f := range filters {
		if f.Condition {
			progressChan <- ProgressMessage{
				Phase:       f.PhaseName,
				Progress:    (filtersApplied * 100) / totalFilters,
				Description: fmt.Sprintf("Starting %s...", f.PhaseName),
			}

			packages = applyConcurrentFilter(packages, f.Filter)
			filtersApplied++

			progressChan <- ProgressMessage{
				Phase:       f.PhaseName,
				Progress:    (filtersApplied * 100) / totalFilters,
				Description: fmt.Sprintf("%s completed", f.PhaseName),
			}
		}
	}

	progressChan <- ProgressMessage{
		Phase:       "Filtering complete",
		Progress:    100,
		Description: "All filtering completed",
	}

	return packages
}
