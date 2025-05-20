package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

func fetchCasks(
	origin string,
	prefix string,
) ([]*pkgdata.PkgInfo, error) {
	caskroomRoot := filepath.Join(prefix, caskroomSubpath)
	installedNames, err := getInstalledCasks(caskroomRoot)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
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

	stage1Out, stage1Err := worker.RunWorkers(
		inputChan,
		func(name string) (*pkgdata.PkgInfo, error) {
			receiptPath := filepath.Join(caskroomRoot, name, ".metadata", receiptName)
			return parseCaskReceipt(name, receiptPath)
		},
		0,
		len(installedNames),
	)

	stage2Out, stage2Err := worker.RunWorkers(
		stage1Out,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			size, err := getInstallSize(filepath.Join(caskroomRoot, pkg.Name, pkg.Version))
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
		return []*pkgdata.PkgInfo{}, metaErr
	}

	stage3Out, stage3Err := worker.RunWorkers(
		stage2Out,
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

	allErrs := worker.MergeErrors(stage1Err, stage2Err, stage3Err)
	return worker.CollectOutput(stage3Out, allErrs)
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
