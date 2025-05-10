package worker

import (
	"errors"
	"fmt"
	"qp/internal/pkgdata"
	"runtime"
	"sync"
)

const DefaultBufferSize = 64

func RunWorkers[T any](
	inputChan <-chan T,
	workerFunc func(T) (*pkgdata.PkgInfo, error),
	numWorkers int, // pass 0 unless testing or intentionally limiting
	bufferSize int,
) ([]*pkgdata.PkgInfo, error) {
	if bufferSize < 1 {
		return nil, fmt.Errorf("invalid buffer size: %d (must be >= 1)", bufferSize)
	}

	if numWorkers <= 0 {
		// fun fact: NumCPU() does account for hyperthreading
		numWorkers = getWorkerCount(runtime.NumCPU())
	}

	outputChan := make(chan *pkgdata.PkgInfo, bufferSize)
	errChan := make(chan error, DefaultBufferSize)

	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range inputChan {
				pkg, err := workerFunc(item)
				if err != nil {
					errChan <- err
					continue
				}

				outputChan <- pkg
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputChan)
		close(errChan)
	}()

	var pkgs []*pkgdata.PkgInfo
	var errs []error

	for pkg := range outputChan {
		pkgs = append(pkgs, pkg)
	}

	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return pkgs, errors.Join(errs...)
	}

	return pkgs, nil
}

func getWorkerCount(numCPUs int) int {
	numWorkers := numCPUs * 2
	if numCPUs <= 2 {
		// let's keep it simple for embedded devices
		numWorkers = numCPUs
	}

	return min(numWorkers, 12) // avoid overthreading on high-core systems
}
