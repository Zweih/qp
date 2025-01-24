package pkgdata

import (
	"sync"
	"time"
)

type FilterCondition struct {
	Condition bool
	Filter    func([]PackageInfo) []PackageInfo
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

func ConcurrentFilters(packages []PackageInfo, dateFilter time.Time, explicitOnly bool, dependenciesOnly bool) []PackageInfo {
	type FilterCondition struct {
		Condition bool
		Filter    func([]PackageInfo) []PackageInfo
	}

	filters := []FilterCondition{
		{
			Condition: explicitOnly,
			Filter:    FilterExplicit,
		},
		{
			Condition: dependenciesOnly,
			Filter:    FilterDependencies,
		},
		{
			Condition: !dateFilter.IsZero(),
			Filter: func(pkgs []PackageInfo) []PackageInfo {
				return FilterByDate(pkgs, dateFilter)
			},
		},
	}

	for _, f := range filters {
		if f.Condition {
			packages = applyConcurrentFilter(packages, f.Filter)
		}
	}

	return packages
}
