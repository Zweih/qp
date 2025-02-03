package pkgdata

import (
	"fmt"
	"sort"
	"sync"
)

const concurrentSortThreshold = 2000

type PackageComparator func(a PackageInfo, b PackageInfo) bool

func alphabeticalComparator(a PackageInfo, b PackageInfo) bool {
	return a.Name < b.Name
}

func dateComparator(a PackageInfo, b PackageInfo) bool {
	return a.Timestamp.Before(b.Timestamp)
}

func sizeDecComparator(a PackageInfo, b PackageInfo) bool {
	return a.Size < b.Size
}

func sizeAscComparator(a PackageInfo, b PackageInfo) bool {
	return a.Size > b.Size
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

func mergedSortedChunks(
	leftChunk []PackageInfo,
	rightChunk []PackageInfo,
	comparator PackageComparator,
) []PackageInfo {
	capacity := len(leftChunk) + len(rightChunk)
	result := make([]PackageInfo, 0, capacity)
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

func sortConcurrently(
	pkgs []PackageInfo,
	comparator PackageComparator,
	phase string,
	reportProgress ProgressReporter,
) []PackageInfo {
	total := len(pkgs)

	if total == 0 {
		return nil
	}

	chunkSize := 100

	var mu sync.Mutex
	var wg sync.WaitGroup

	numChunks := (total + chunkSize - 1) / chunkSize
	chunks := make([][]PackageInfo, 0, numChunks) // pre-allocate

	for chunkIdx := 0; chunkIdx < numChunks; chunkIdx++ {
		startIdx := chunkIdx * chunkSize
		endIdx := startIdx + chunkSize

		if endIdx > total {
			endIdx = total
		}

		chunk := make([]PackageInfo, endIdx-startIdx)
		copy(chunk, pkgs[startIdx:endIdx]) // avoid mutating the original array

		wg.Add(1)

		// TODO: c is a bad name, maybe make this a seperate function
		go func(c []PackageInfo) {
			defer wg.Done()

			sort.SliceStable(c, func(i int, j int) bool {
				return comparator(c[i], c[j])
			})

			mu.Lock()
			chunks = append(chunks, c)
			mu.Unlock()

			if reportProgress != nil {
				currentProgress := (chunkIdx + 1) * 50 / numChunks // scale chunk sorting progress to 0%-50%
				reportProgress(currentProgress, 100, fmt.Sprintf("%s - Sorted chunk %d/%d", phase, chunkIdx+1, numChunks))
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
		var newChunks [][]PackageInfo

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

func sortNormally(
	pkgs []PackageInfo,
	comparator PackageComparator,
	phase string,
	reportProgress ProgressReporter,
) []PackageInfo {
	sortedPkgs := make([]PackageInfo, len(pkgs))
	copy(sortedPkgs, pkgs)

	reportProgress(0, 100, fmt.Sprintf("%s - normally", phase))

	sort.SliceStable(sortedPkgs, func(i int, j int) bool {
		return comparator(sortedPkgs[i], sortedPkgs[j])
	})

	reportProgress(100, 100, fmt.Sprintf("%s completed", phase))

	return sortedPkgs
}

func SortPackages(pkgs []PackageInfo, sortBy string, reportProgress ProgressReporter) []PackageInfo {
	comparator, valid := getComparator(sortBy)

	if !valid {
		panic(fmt.Sprintf("invalid sort mode: %s", sortBy))
	}

	phase := "Sorting packages"

	if len(pkgs) < concurrentSortThreshold {
		return sortNormally(pkgs, comparator, phase, reportProgress)
	}

	return sortConcurrently(pkgs, comparator, phase, reportProgress)
}
