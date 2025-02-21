package pkgdata

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	fieldName        = "%NAME%"
	fieldInstallDate = "%INSTALLDATE%"
	fieldSize        = "%SIZE%"
	fieldReason      = "%REASON%"
	fieldVersion     = "%VERSION%"
	pacmanDbPath     = "/var/lib/pacman/local"
)

func FetchPackages() ([]PackageInfo, error) {
	packagePaths, err := os.ReadDir(pacmanDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pacman database: %v", err)
	}

	numPackages := len(packagePaths)

	var wg sync.WaitGroup
	descPaths := make(chan string, numPackages)
	packagesChan := make(chan PackageInfo, numPackages)
	errorsChan := make(chan error, numPackages)

	// fun fact: NumCPU() does account for hyperthreading
	numWorkers := getWorkerCount(runtime.NumCPU(), numPackages)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for descPath := range descPaths {
				pkg, err := parseDescFile(descPath)
				if err != nil {
					errorsChan <- err
					continue
				}

				packagesChan <- pkg
			}
		}()
	}

	for _, packagePath := range packagePaths {
		if packagePath.IsDir() {
			descPath := filepath.Join(pacmanDbPath, packagePath.Name(), "desc")
			descPaths <- descPath
		}
	}

	close(descPaths)

	wg.Wait()
	close(packagesChan)
	close(errorsChan)

	packages := make([]PackageInfo, 0, numPackages)
	for pkg := range packagesChan {
		packages = append(packages, pkg)
	}

	if len(errorsChan) > 0 {
		var collectedErrors []error

		for err := range errorsChan {
			collectedErrors = append(collectedErrors, err)
		}

		return nil, errors.Join(collectedErrors...)
	}

	return packages, nil
}

func getWorkerCount(numCPUs int, numFiles int) int {
	var numWorkers int

	if numCPUs <= 2 {
		// let's keep it simple for devices like rPi zeroes
		numWorkers = numCPUs
	} else {
		numWorkers = numCPUs * 2
	}

	if numWorkers > numFiles {
		return numFiles // don't use more workers than files
	}

	return min(numWorkers, 12) // avoid overthreading on high-core systems
}

func parseDescFile(descPath string) (PackageInfo, error) {
	file, err := os.Open(descPath)
	if err != nil {
		return PackageInfo{}, fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()

	var pkg PackageInfo
	var currentField string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch line {
		case fieldName, fieldInstallDate, fieldSize, fieldReason, fieldVersion:
			currentField = line
		default:
			if err := applyField(&pkg, currentField, line); err != nil {
				return PackageInfo{}, fmt.Errorf("error reading desc file %s: %w", descPath, err)
			}

			currentField = "" // reset
		}

	}

	if pkg.Name == "" {
		return PackageInfo{}, fmt.Errorf("package name is missing in file: %s", descPath)
	}

	if pkg.Reason == "" {
		pkg.Reason = "explicit"
	}

	return pkg, nil
}

func applyField(pkg *PackageInfo, field string, value string) error {
	switch field {
	case fieldName:
		pkg.Name = value

	case fieldReason:
		if value == "1" {
			pkg.Reason = "dependency"
		} else {
			pkg.Reason = "explicit"
		}

	case fieldInstallDate:
		installDate, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid install date value %q: %w", value, err)
		}

		pkg.Timestamp = time.Unix(installDate, 0)

	case fieldVersion:
		pkg.Version = value

	case fieldSize:
		size, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid size value %q: %w", value, err)
		}

		pkg.Size = size

	default:
		// ignore unknown fields
	}

	return nil
}
