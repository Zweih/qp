package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	pb "qp/internal/protobuf"
	"qp/internal/worker"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

const pkgsPerChunk = 100

func SaveCacheModTime(cacheRoot string, modTime int64) error {
	modPath := cacheRoot + dotModTime
	return FileManager.WriteFile(modPath, []byte(strconv.FormatInt(modTime, 10)), 0644)
}

func LoadCacheModTime(cacheRoot string) (int64, error) {
	modPath := cacheRoot + dotModTime
	data, err := os.ReadFile(modPath)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

func SaveProtoCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	cleanupOldChunks(cacheRoot)

	for i := 0; i < len(pkgs); i += pkgsPerChunk {
		end := min(i+pkgsPerChunk, len(pkgs))
		chunk := pkgs[i:end]

		cachedChunk := &pb.CachedPkgs{
			Pkgs:    pkgsToProtos(chunk),
			Version: cacheVersion,
		}

		chunkCount := i / pkgsPerChunk

		byteData, err := proto.Marshal(cachedChunk)
		if err != nil {
			return fmt.Errorf("failed to marshal chunk %d: %v", chunkCount, err)
		}

		chunkPath := cacheRoot + "." + strconv.Itoa(chunkCount) + dotCache
		err = FileManager.WriteFile(chunkPath, byteData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write chunk %d: %v", chunkCount, err)
		}
	}

	return nil
}

func LoadProtoCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	chunkPaths := findChunkFiles(cacheRoot)
	if len(chunkPaths) == 0 {
		return nil, errors.New("no cache chunks found")
	}

	if len(chunkPaths) == 1 {
		return loadSingleChunk(chunkPaths[0])
	}

	inputChan := make(chan string, len(chunkPaths))
	for _, path := range chunkPaths {
		inputChan <- path
	}
	close(inputChan)

	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	resultChan := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		loadSingleChunk,
		runtime.NumCPU(),
		len(chunkPaths),
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	chunks, err := worker.CollectOutput(resultChan, errChan)
	if err != nil {
		return nil, err
	}

	allPkgs := make([]*pkgdata.PkgInfo, 0, len(chunkPaths)*pkgsPerChunk)
	for _, chunk := range chunks {
		allPkgs = append(allPkgs, chunk...)
	}

	return allPkgs, nil
}

func loadSingleChunk(chunkPath string) ([]*pkgdata.PkgInfo, error) {
	data, err := os.ReadFile(chunkPath)
	if err != nil {
		return nil, err
	}

	var chunk pb.CachedPkgs
	err = proto.Unmarshal(data, &chunk)
	if err != nil {
		return nil, err
	}

	if chunk.Version != cacheVersion {
		return nil, errors.New("chunk version mismatch")
	}

	return protosToPkgs(chunk.Pkgs), nil
}

func findChunkFiles(cacheRoot string) []string {
	pattern := cacheRoot + ".[0-9]*" + dotCache
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	return matches
}

func cleanupOldChunks(cacheRoot string) {
	pattern := cacheRoot + "*" + dotCache
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	for _, path := range matches {
		os.Remove(path)
	}
}
