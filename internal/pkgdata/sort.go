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
	chunkSize := min(100, baseChunkSize)

	var mu sync.Mutex
	var wg sync.WaitGroup

	numChunks := (total + chunkSize - 1) / chunkSize
	chunks := make([][]*PkgInfo, 0, numChunks) // pre-allocate

	for chunkIdx := range numChunks {
		startIdx := chunkIdx * chunkSize
		endIdx := min(startIdx+chunkSize, total)

		chunk := pkgPtrs[startIdx:endIdx]

		wg.Add(1)

		go func(c []*PkgInfo) {
			defer wg.Done()

			sort.Slice(c, func(i int, j int) bool {
				return comparator(c[i], c[j])
			})

			mu.Lock()
			chunks = append(chunks, c)
			mu.Unlock()

			if reportProgress != nil {
				currentProgress := (chunkIdx + 1) * 50 / numChunks // scale chunk sorting progress to 0%-50%
				reportProgress(
					currentProgress,
					100,
					fmt.Sprintf("%s - Sorted chunk %d/%d", phase, chunkIdx+1, numChunks),
				)
			}
		}(chunk)
	}

	wg.Wait()

	if reportProgress != nil {
		// "halfway" there
		reportProgress(50, 100, fmt.Sprintf("%s - Initial chunk sorting complete", phase))
	}

	mergeStep := 0

	for len(chunks) > 1 {
		var newChunks [][]*PkgInfo

		for i := 0; i < len(chunks); i += 2 {
			if i+1 < len(chunks) {
				mergedChunk := mergedSortedChunks(chunks[i], chunks[i+1], comparator)
				newChunks = append(newChunks, mergedChunk)

				continue
			}

			newChunks = append(newChunks, chunks[i]) // move odd chunk forward
		}

		chunks = newChunks

		if reportProgress != nil {
			mergeStep++
			currentProgress := 50 + (mergeStep * 50 / (numChunks - 1)) // scale to 50%-100%
			reportProgress(currentProgress, 100, fmt.Sprintf("%s - Merging step %d", phase, mergeStep))
		}
	}

	if reportProgress != nil {
		reportProgress(total, total, fmt.Sprintf("%s completed", phase))
	}

	if len(chunks) == 1 {
		return chunks[0]
	}

	return nil
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

	sort.Slice(pkgPtrs, func(i int, j int) bool {
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
