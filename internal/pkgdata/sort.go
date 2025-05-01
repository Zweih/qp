package pkgdata

import (
	"errors"
	"fmt"
	"qp/internal/consts"
	"qp/internal/pipeline/meta"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const ConcurrentSortThreshold = 500

type PkgComparator func(a *PkgInfo, b *PkgInfo) bool

type ordered interface {
	~int64 | ~string
}

func makeComparator[T ordered](
	getValue func(*PkgInfo) T,
	asc bool,
) PkgComparator {
	if asc {
		return func(a, b *PkgInfo) bool { return getValue(a) < getValue(b) }
	}

	return func(a, b *PkgInfo) bool { return getValue(a) > getValue(b) }
}

func GetComparator(field consts.FieldType, asc bool) (PkgComparator, error) {
	switch field {
	case consts.FieldDate, consts.FieldBuildDate, consts.FieldSize:
		return makeComparator(func(p *PkgInfo) int64 { return p.GetInt(field) }, asc), nil

	case consts.FieldName, consts.FieldLicense,
		consts.FieldPkgType, consts.FieldPkgBase,
		consts.FieldPackager, consts.FieldArch:
		return makeComparator(func(p *PkgInfo) string {
			return strings.ToLower(p.GetString(field))
		}, asc), nil

	default:
		return nil, errors.New("invalid sort field")
	}
}

func mergedSortedChunks(
	leftChunk []*PkgInfo,
	rightChunk []*PkgInfo,
	comparator PkgComparator,
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
func SortConcurrently(
	pkgPtrs []*PkgInfo,
	comparator PkgComparator,
	phase string,
	reportProgress meta.ProgressReporter,
) []*PkgInfo {
	total := len(pkgPtrs)

	if total == 0 {
		return nil
	}

	numCPUs := runtime.NumCPU()
	baseChunkSize := total / (2 * numCPUs)
	chunkSize := max(100, baseChunkSize)

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
func SortNormally(
	pkgPtrs []*PkgInfo,
	comparator PkgComparator,
	phase string,
	reportProgress meta.ProgressReporter,
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
