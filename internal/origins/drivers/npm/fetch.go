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

func fetchPackages(
	origin string,
	modulesDir string,
	outChan chan<- *pkgdata.PkgInfo,
	errChan chan<- error,
	errGroup *sync.WaitGroup,
) {
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		errChan <- fmt.Errorf("failed to read global node_modules directory: %w", err)
		return
	}

	nodeVersion := extractNodeVersion(modulesDir)
	inputChan := make(chan string, len(entries))

	var packagePaths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		if entryName != "" && entryName[0] == '@' {
			scope := filepath.Join(modulesDir, entryName)
			subEntries, err := os.ReadDir(scope)
			if err != nil {
				errChan <- fmt.Errorf("failed to read scope directory %s: %w", scope, err)
				continue
			}

			for _, subEntry := range subEntries {
				if !subEntry.IsDir() {
					continue
				}

				inputChan <- filepath.Join(entryName, subEntry.Name())
			}

			continue
		}

		inputChan <- entryName
	}

	if len(packagePaths) == 0 {
		return
	}

	close(inputChan)

	stage1 := worker.RunWorkers(
		inputChan,
		errChan,
		errGroup,
		func(pkgName string) (*pkgdata.PkgInfo, error) {
			return parsePackageJson(filepath.Join(modulesDir, pkgName))
		},
		0,
		len(packagePaths),
	)

	stage2 := worker.RunWorkers(
		stage1,
		errChan,
		errGroup,
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

			pkg.Name = name
			pkg.Groups = groups
			pkg.Size = size
			pkg.Origin = origin
			pkg.Env = nodeVersion

			return pkg, nil
		},
		0,
		len(packagePaths),
	)

	for pkg := range stage2 {
		outChan <- pkg
	}
}

func extractNodeVersion(modulesDir string) string {
	if strings.Contains(modulesDir, nvmDir) {
		parts := strings.Split(modulesDir, "/")
		for i, part := range parts {
			if part == "node" && i+1 < len(parts) {
				return nvmNodeEnv + parts[i+1]
			}
		}
	}

	if strings.HasPrefix(modulesDir, "/usr/") {
		return nodeSystemEnv
	}

	return nodeUnknownEnv
}
