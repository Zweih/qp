package npm

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/origins/shared"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"
	"sync"
)

func fetchPackages(origin string, modulesDir string) ([]*pkgdata.PkgInfo, error) {
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read global node_modules directory: %w", err)
	}

	inputChan := make(chan string, len(entries))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		if entryName != "" && entryName[0] == '@' {
			scope := filepath.Join(modulesDir, entryName)
			if subEntries, err := os.ReadDir(scope); err == nil {
				for _, subEntry := range subEntries {
					if !subEntry.IsDir() {
						continue
					}

					inputChan <- filepath.Join(entryName, subEntry.Name())
				}
			}

			continue
		}

		inputChan <- entryName
	}

	close(inputChan)

	stage1 := worker.RunWorkers(
		inputChan,
		errChan,
		&errGroup,
		func(pkgName string) (*pkgdata.PkgInfo, error) {
			return parsePackageJson(filepath.Join(modulesDir, pkgName))
		},
		0,
		len(entries),
	)

	stage2 := worker.RunWorkers(
		stage1,
		errChan,
		&errGroup,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			pkgDir := filepath.Join(modulesDir, pkg.Name)
			size, err := shared.GetInstallSize(pkgDir)
			if err != nil {
				return nil, err
			}

			var groups []string
			name := pkg.Name

			if len(name) > 1 && name[0] == '@' {
				parts := strings.SplitN(name, "/", 2)
				if len(parts) == 2 {
					groups = []string{parts[0]}
					name = parts[1]
				}
			}

			creationTime, isReliable, err := shared.GetCreationTime(filepath.Dir(pkgDir))
			if err == nil && isReliable {
				pkg.InstallTimestamp = creationTime
			}

			pkg.Name = name
			pkg.Groups = groups
			pkg.Size = size
			pkg.Origin = origin

			return pkg, nil
		},
		0,
		len(entries),
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	return worker.CollectOutput(stage2, errChan)
}
