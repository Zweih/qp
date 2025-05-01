package pacman

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	"runtime"
	"sync"
)

func fetchPackages(origin string) ([]*pkgdata.PkgInfo, error) {
	pkgPaths, err := os.ReadDir(PacmanDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pacman database: %v", err)
	}

	numPkgs := len(pkgPaths)

	var wg sync.WaitGroup
	descPathChan := make(chan string, numPkgs)
	pkgChan := make(chan *pkgdata.PkgInfo, numPkgs)
	errorsChan := make(chan error, numPkgs)

	// fun fact: NumCPU() does account for hyperthreading
	numWorkers := getWorkerCount(runtime.NumCPU(), numPkgs)

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for descPath := range descPathChan {
				pkg, err := parseDescFile(descPath)
				if err != nil {
					errorsChan <- err
					continue
				}

				pkg.Origin = origin
				pkgChan <- pkg
			}
		}()
	}

	for _, packagePath := range pkgPaths {
		if packagePath.IsDir() {
			descPath := filepath.Join(PacmanDbPath, packagePath.Name(), "desc")
			descPathChan <- descPath
		}
	}

	close(descPathChan)

	wg.Wait()
	close(pkgChan)
	close(errorsChan)

	if len(errorsChan) > 0 {
		var collectedErrors []error

		for err := range errorsChan {
			collectedErrors = append(collectedErrors, err)
		}

		return nil, errors.Join(collectedErrors...)
	}

	pkgs := make([]*pkgdata.PkgInfo, 0, numPkgs)
	for pkg := range pkgChan {
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func getWorkerCount(numCPUs int, numFiles int) int {
	var numWorkers int

	numWorkers = numCPUs * 2
	if numCPUs <= 2 {
		// let's keep it simple for devices like rPi zeroes
		numWorkers = numCPUs
	}

	if numWorkers > numFiles {
		return numFiles // don't use more workers than files
	}

	return min(numWorkers, 12) // avoid overthreading on high-core systems
}
