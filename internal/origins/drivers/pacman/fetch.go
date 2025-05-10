package pacman

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
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

	outputChan, errChan := worker.RunWorkers(
		descPathChan,
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

	return worker.CollectOutput(outputChan, errChan)
}
