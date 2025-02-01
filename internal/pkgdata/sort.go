package pkgdata

import (
	"fmt"
	"sort"
)

type PackageComparator func(pkgs []PackageInfo, i int, j int) bool

func alphabeticalComparator(pkgs []PackageInfo, i int, j int) bool {
	return pkgs[i].Name < pkgs[j].Name
}

func dateComparator(pkgs []PackageInfo, i int, j int) bool {
	return pkgs[i].Timestamp.Before(pkgs[j].Timestamp)
}

func sizeDecComparator(pkgs []PackageInfo, i int, j int) bool {
	return pkgs[i].Size < pkgs[j].Size
}

func sizeAscComparator(pkgs []PackageInfo, i int, j int) bool {
	return pkgs[i].Size > pkgs[j].Size
}

func getComparator(sortBy string) (PackageComparator, bool) {
	switch sortBy {
	case "alphabetical":
		return alphabeticalComparator, true
	case "date":
		return dateComparator, true
	case "size:dec":
		return sizeDecComparator, true
	case "size:asc":
		return sizeAscComparator, true
	default:
		return nil, false
	}
}

func sortWithProgress(pkgs []PackageInfo, comparator PackageComparator, phase string, reportProgress ProgressReporter) {
	total := len(pkgs)
	progressStep := max(1, total/10) // wait they finally added min() and max() to standard lib??

	sort.SliceStable(pkgs, func(i int, j int) bool {
		remainder := i % progressStep

		if reportProgress != nil && remainder == 0 {
			reportProgress(i, total, phase)
		}

		return comparator(pkgs, i, j)
	})

	if reportProgress != nil {
		reportProgress(total, total, fmt.Sprintf("%s completed", phase))
	}
}

func SortPackages(pkgs []PackageInfo, sortBy string, reportProgress ProgressReporter) []PackageInfo {
	sortedPkgs := make([]PackageInfo, len(pkgs))
	copy(sortedPkgs, pkgs)

	comparator, valid := getComparator(sortBy)

	if !valid {
		panic(fmt.Sprintf("invalid sort mode: %s", sortBy))
	}

	sortWithProgress(sortedPkgs, comparator, "Sorting packages", reportProgress)

	return sortedPkgs
}
