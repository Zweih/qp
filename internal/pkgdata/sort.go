package pkgdata

import (
	"errors"
	"qp/internal/consts"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const ConcurrentSortThreshold = 500

type PkgComparator func(a *PkgInfo, b *PkgInfo) bool

type ordered interface {
	~int64 | ~string | ~int
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
	switch consts.GetFieldPrim(field) {
	case consts.FieldPrimDate, consts.FieldPrimSize:
		return makeComparator(func(p *PkgInfo) int64 { return p.GetInt(field) }, asc), nil

	case consts.FieldPrimStr:
		return makeComparator(func(p *PkgInfo) string {
			return strings.ToLower(p.GetString(field))
		}, asc), nil

	case consts.FieldPrimStrArr:
		return makeComparator(func(p *PkgInfo) int {
			return len(p.GetStrArr(field))
		}, asc), nil

	case consts.FieldPrimRel:
		return makeComparator(func(p *PkgInfo) int {
			return len(GetRelationsByDepth(p.GetRelations(field), 1))
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
	pkg []*PkgInfo,
	comparator PkgComparator,
) []*PkgInfo {
	total := len(pkg)

	if total == 0 {
		return nil
	}

	numCPUs := runtime.NumCPU()
	baseChunkSize := total / (2 * numCPUs)
	chunkSize := max(100, baseChunkSize)

	var mu sync.Mutex
	var wg sync.WaitGroup

	numChunks := (total + chunkSize - 1) / chunkSize
	chunks := make([][]*PkgInfo, 0, numChunks)

	for chunkIdx := range numChunks {
		startIdx := chunkIdx * chunkSize
		endIdx := min(startIdx+chunkSize, total)

		chunk := pkg[startIdx:endIdx]

		wg.Add(1)

		go func(c []*PkgInfo) {
			defer wg.Done()

			sort.Slice(c, func(i int, j int) bool {
				return comparator(c[i], c[j])
			})

			mu.Lock()
			chunks = append(chunks, c)
			mu.Unlock()
		}(chunk)
	}

	wg.Wait()

	for len(chunks) > 1 {
		var newChunks [][]*PkgInfo

		for i := 0; i < len(chunks); i += 2 {
			if i+1 < len(chunks) {
				mergedChunk := mergedSortedChunks(chunks[i], chunks[i+1], comparator)
				newChunks = append(newChunks, mergedChunk)

				continue
			}

			newChunks = append(newChunks, chunks[i])
		}

		chunks = newChunks
	}

	if len(chunks) == 1 {
		return chunks[0]
	}

	return nil
}

// pkgPointers will be sorted in place, mutating the slice order
func SortNormally(
	pkgs []*PkgInfo,
	comparator PkgComparator,
) []*PkgInfo {
	sort.Slice(pkgs, func(i int, j int) bool {
		return comparator(pkgs[i], pkgs[j])
	})

	return pkgs
}
