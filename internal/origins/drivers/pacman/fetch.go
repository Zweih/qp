package pacman

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

func fetchPackages(origin string) ([]*pkgdata.PkgInfo, error) {
	pkgPaths, err := os.ReadDir(PacmanDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pacman database: %v", err)
	}

	numPkgs := len(pkgPaths)
	descPathChan := make(chan string, numPkgs)

	go func() {
		for _, packagePath := range pkgPaths {
			if packagePath.IsDir() {
				descPath := filepath.Join(PacmanDbPath, packagePath.Name(), "desc")
				descPathChan <- descPath
			}
		}

		close(descPathChan)
	}()

	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	outputChan := worker.RunWorkers(
		descPathChan,
		errChan,
		&errGroup,
		func(path string) (*pkgdata.PkgInfo, error) {
			pkg, err := parseDescFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", path, err)
			}

			pkg.Origin = origin

			return pkg, nil
		},
		0,
		numPkgs,
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(outputChan, errChan)
}
