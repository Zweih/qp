package pkgdata

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"yaylog/internal/config"
)

const concurrentSortThreshold = 500

type PackageComparator func(a *PkgInfo, b *PkgInfo) bool

func alphabeticalComparator(a *PkgInfo, b *PkgInfo) bool {
	return a.Name < b.Name
}

func dateComparator(a *PkgInfo, b *PkgInfo) bool {
	return a.Timestamp < b.Timestamp
}

func sizeDecComparator(a *PkgInfo, b *PkgInfo) bool {
	return a.Size > b.Size
}

func sizeAscComparator(a *PkgInfo, b *PkgInfo) bool {
	return a.Size < b.Size
}

func getComparator(sortBy string) PackageComparator {
	switch sortBy {
	case "alphabetical":
		return alphabeticalComparator
	case "date":
		return dateComparator
	case "size:desc":
		return sizeDecComparator
	case "size:asc":
		return sizeAscComparator
	default:
		return nil
	}
}

func mergedSortedChunks(
	leftChunk []*PkgInfo,
	rightChunk []*PkgInfo,
	comparator PackageComparator,
) []*PkgInfo {
	capacity := len(leftChunk) + len(rightChunk)
	result := make([]*PkgInfo, 0, capacity)
	i, j := 0, 0

	for i < len(leftChunk) && j < len(rightChunk) {
		if comparator(leftChunk[i], rightChunk[j]) {
			result = append(result, leftChunk[i])
			i++
			continue
		}

		result = append(result, rightChunk[j])
		j++
	}

	// append remaining elements
	result = append(result, leftChunk[i:]...)
	result = append(result, rightChunk[j:]...)

	return result
}

// pkgPointers will be sorted in place, mutating the slice order
func sortConcurrently(
	pkgPtrs []*PkgInfo,
	comparator PackageComparator,
	phase string,
	reportProgress ProgressReporter,
) []*PkgInfo {
	total := len(pkgPtrs)

	if total == 0 {
		return nil
	}

	numCPUs := runtime.NumCPU()
	baseChunkSize := total / (2 * numCPUs)
	chunkSize := max(100, baseChunkSize)

	var wg sync.WaitGroup
	numChunks := (total + chunkSize - 1) / chunkSize

	for chunkIdx := range numChunks {
		startIdx := chunkIdx * chunkSize
		endIdx := min(startIdx+chunkSize, total)

		wg.Add(1)

		go func() {
			defer wg.Done()

			sort.SliceStable(pkgPtrs[startIdx:endIdx], func(i int, j int) bool {
				return comparator(pkgPtrs[i], pkgPtrs[j])
			})

			if reportProgress != nil {
				currentProgress := (chunkIdx + 1) * 50 / numChunks // scale chunk sorting progress to 0%-50%
				reportProgress(
					currentProgress,
					100,
					fmt.Sprintf("%s - Sorted chunk %d/%d", phase, chunkIdx+1, numChunks),
				)
			}
		}()
	}

	wg.Wait()

	if reportProgress != nil {
		// "halfway" there
		reportProgress(50, 100, fmt.Sprintf("%s - Initial chunk sorting complete", phase))
	}

	sort.SliceStable(pkgPtrs, func(i int, j int) bool {
		return comparator(pkgPtrs[i], pkgPtrs[j])
	})

	if reportProgress != nil {
		reportProgress(total, total, fmt.Sprintf("%s completed", phase))
	}

	return pkgPtrs
}

// pkgPointers will be sorted in place, mutating the slice order
func sortNormally(
	pkgPtrs []*PkgInfo,
	comparator PackageComparator,
	phase string,
	reportProgress ProgressReporter,
) []*PkgInfo {
	if reportProgress != nil {
		reportProgress(0, 100, fmt.Sprintf("%s - normally", phase))
	}

	sort.SliceStable(pkgPtrs, func(i int, j int) bool {
		return comparator(pkgPtrs[i], pkgPtrs[j])
	})

	if reportProgress != nil {
		reportProgress(100, 100, fmt.Sprintf("%s completed", phase))
	}

	return pkgPtrs
}

func SortPackages(
	cfg config.Config,
	pkgPtrs []*PkgInfo,
	reportProgress ProgressReporter,
) ([]*PkgInfo, error) {
	comparator := getComparator(cfg.SortBy)
	phase := "Sorting packages"

	// threshold is 500 as that is where merge sorting chunk performance overtakes timsort
	if len(pkgPtrs) < concurrentSortThreshold {
		return sortNormally(pkgPtrs, comparator, phase, reportProgress), nil
	}

	return sortConcurrently(pkgPtrs, comparator, phase, reportProgress), nil
}
