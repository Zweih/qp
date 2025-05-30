package pacman

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"qp/internal/origins/shared"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"sync"
)

func fetchPackages(origin string, cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	pkgPaths, err := os.ReadDir(pacmanDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pacman database: %v", err)
	}

	numPkgs := len(pkgPaths)
	descPathChan := make(chan string, numPkgs)

	cachedHistory, latestLogTime, err := pkgdata.LoadInstallHistory(cacheRoot)
	if err != nil {
		cachedHistory = make(map[string]int64)
		latestLogTime = 0
	}

	freshHistory, newLatestTime, err := parseLogHistory(latestLogTime)
	if err != nil {
		freshHistory = make(map[string]int64)
		newLatestTime = latestLogTime
	}

	combinedHistory := make(map[string]int64)
	currentHistory := make(map[string]int64)

	for name, timestamp := range cachedHistory {
		combinedHistory[name] = timestamp
	}

	for name, timestamp := range freshHistory {
		combinedHistory[name] = timestamp
	}

	systemInstallTime, err := getSystemInstallTime()
	if err != nil {
		return nil, err
	}

	go func() {
		for _, packagePath := range pkgPaths {
			if packagePath.IsDir() {
				descPath := filepath.Join(pacmanDbPath, packagePath.Name(), "desc")
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
			installTime := systemInstallTime

			if logInstallTime, exists := combinedHistory[pkg.Name]; exists {
				installTime = logInstallTime
			}

			pkg.InstallTimestamp = installTime
			return pkg, nil
		},
		0,
		numPkgs,
	)

	go func() {
		errGroup.Wait()
		close(errChan)
	}()

	pkgs, err := worker.CollectOutput(outputChan, errChan)
	if err != nil {
		return pkgs, err
	}

	for _, pkg := range pkgs {
		if installTime, exists := combinedHistory[pkg.Name]; exists {
			currentHistory[pkg.Name] = installTime
		}
	}

	err = pkgdata.SaveInstallHistory(cacheRoot, currentHistory, newLatestTime)
	return pkgs, err
}

func getSystemInstallTime() (int64, error) {
	systemPaths := []string{
		"/etc/machine-id",
		"/etc/hostname",
		"/boot/vmlinuz-linux",
		"/usr/bin/pacman",
		"/var/lib/pacman",
		"/etc/passwd",
		"/etc/group",
	}

	var oldestTime int64 = math.MaxInt64
	foundAny := false

	for _, path := range systemPaths {
		if fileInfo, err := os.Stat(path); err == nil {
			if birthTime, reliable, err := shared.GetBirthTime(path); err == nil && reliable {
				if birthTime < oldestTime {
					oldestTime = birthTime
					foundAny = true
				}

				continue
			}

			modTime := fileInfo.ModTime().Unix()
			if modTime < oldestTime {
				oldestTime = modTime
				foundAny = true
			}
		}
	}

	if !foundAny {
		return 0, fmt.Errorf("could not determine system install time")
	}

	return oldestTime, nil
}
