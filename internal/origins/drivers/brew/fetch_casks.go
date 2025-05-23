package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/origins/shared"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

func fetchCasks(
	origin string,
	prefix string,
	outChan chan<- *pkgdata.PkgInfo,
	errChan chan<- error,
	errGroup *sync.WaitGroup,
) {
	caskroomRoot := filepath.Join(prefix, caskroomSubpath)
	installedNames, err := getInstalledCasks(caskroomRoot)
	if err != nil {
		errChan <- err
		return
	}

	if len(installedNames) < 1 {
		return
	}

	wanted := make(map[string]struct{}, len(installedNames))
	for _, name := range installedNames {
		wanted[name] = struct{}{}
	}

	var caskMeta map[string]*CaskMetadata
	var metaErr error
	var metaWg sync.WaitGroup

	metaWg.Add(1)
	go func() {
		defer metaWg.Done()
		caskMeta, metaErr = loadMetadata(caskCachePath, getCaskKey, wanted)
	}()

	inputChan := make(chan string, len(installedNames))
	for _, name := range installedNames {
		inputChan <- name
	}

	close(inputChan)

	stage1Out := worker.RunWorkers(
		inputChan,
		errChan,
		errGroup,
		func(name string) (*pkgdata.PkgInfo, error) {
			receiptPath := filepath.Join(caskroomRoot, name, ".metadata", receiptName)
			return parseCaskReceipt(name, receiptPath)
		},
		0,
		len(installedNames),
	)

	stage2Out := worker.RunWorkers(
		stage1Out,
		errChan,
		errGroup,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			size, err := shared.GetInstallSize(filepath.Join(caskroomRoot, pkg.Name, pkg.Version))
			if err == nil {
				pkg.Size = size
			}

			return pkg, nil
		},
		0,
		len(installedNames),
	)

	metaWg.Wait()
	if metaErr != nil {
		errChan <- metaErr
		return
	}

	stage3Out := worker.RunWorkers(
		stage2Out,
		errChan,
		errGroup,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			if meta, ok := caskMeta[pkg.Name]; ok {
				mergeCaskMetadata(pkg, meta)
			}

			pkg.Origin = origin

			return pkg, nil
		},
		0,
		len(installedNames),
	)

	for pkg := range stage3Out {
		outChan <- pkg
	}
}

func getInstalledCasks(caskroomRoot string) ([]string, error) {
	entries, err := os.ReadDir(caskroomRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read Caskroom directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		names = append(names, entry.Name())
	}

	return names, nil
}
