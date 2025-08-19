package pkgtool

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	"qp/internal/worker"
	"sync"
)

func fetchPackages(origin string) ([]*pkgdata.PkgInfo, error) {
	entries, err := os.ReadDir(packagesDbPath)
	if err != nil {
		return []*pkgdata.PkgInfo{}, fmt.Errorf("failed to read package install directories: %w", err)
	}

	inputChan := make(chan string, len(entries))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, entry := range entries {
		if !entry.IsDir() {
			inputChan <- filepath.Join(packagesDbPath, entry.Name())
		}
	}

	close(inputChan)

	resultChan := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(packagePath string) (*pkgdata.PkgInfo, error) {
			return parsePackageFile(packagePath, origin)
		},
		0,
		len(entries),
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(resultChan, errChan)
}
