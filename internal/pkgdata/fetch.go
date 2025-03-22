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
)

const (
	fieldName        = "%NAME%"
	fieldInstallDate = "%INSTALLDATE%"
	fieldSize        = "%SIZE%"
	fieldReason      = "%REASON%"
	fieldVersion     = "%VERSION%"
	fieldDepends     = "%DEPENDS%"
	fieldProvides    = "%PROVIDES%"
	fieldConflicts   = "%CONFLICTS%"
	fieldArch        = "%ARCH%"
	fieldLicense     = "%LICENSE%"
	fieldUrl         = "%URL%"

	pacmanDbPath = "/var/lib/pacman/local"
)

func FetchPackages() ([]*PkgInfo, error) {
	pkgPaths, err := os.ReadDir(pacmanDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pacman database: %v", err)
	}

	numPkgs := len(pkgPaths)

	var wg sync.WaitGroup
	descPathChan := make(chan string, numPkgs)
	pkgPtrsChan := make(chan *PkgInfo, numPkgs)
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

				pkgPtrsChan <- pkg
			}
		}()
	}

	for _, packagePath := range pkgPaths {
		if packagePath.IsDir() {
			descPath := filepath.Join(pacmanDbPath, packagePath.Name(), "desc")
			descPathChan <- descPath
		}
	}

	close(descPathChan)

	wg.Wait()
	close(pkgPtrsChan)
	close(errorsChan)

	if len(errorsChan) > 0 {
		var collectedErrors []error

		for err := range errorsChan {
			collectedErrors = append(collectedErrors, err)
		}

		return nil, errors.Join(collectedErrors...)
	}

	pkgPtrs := make([]*PkgInfo, 0, numPkgs)
	for pkg := range pkgPtrsChan {
		pkgPtrs = append(pkgPtrs, pkg)
	}

	return pkgPtrs, nil
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

func parseDescFile(descPath string) (*PkgInfo, error) {
	file, err := os.Open(descPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()

	var pkg PkgInfo
	var currentField string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch line {
		case fieldName,
			fieldInstallDate,
			fieldSize,
			fieldReason,
			fieldVersion,
			fieldDepends,
			fieldProvides,
			fieldConflicts,
			fieldArch,
			fieldLicense,
			fieldUrl:
			currentField = line
		case "":
			currentField = "" // reset if line is blank
		default:
			if err := applyField(&pkg, currentField, line); err != nil {
				return nil, fmt.Errorf("error reading desc file %s: %w", descPath, err)
			}
		}
	}

	if pkg.Name == "" {
		return nil, fmt.Errorf("package name is missing in file: %s", descPath)
	}

	if pkg.Reason == "" {
		pkg.Reason = "explicit"
	}

	return &pkg, nil
}

func applyField(pkg *PkgInfo, field string, value string) error {
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

		pkg.Timestamp = installDate

	case fieldVersion:
		pkg.Version = value

	case fieldSize:
		size, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid size value %q: %w", value, err)
		}

		pkg.Size = size

	case fieldDepends:
		// use this if we ever need to separate the package name from its dependencies re := regexp.MustCompile(`^([^<>=]+)`)
		pkg.Depends = append(pkg.Depends, value)

	case fieldProvides:
		pkg.Provides = append(pkg.Provides, value)

	case fieldConflicts:
		pkg.Conflicts = append(pkg.Conflicts, value)

	case fieldArch:
		pkg.Arch = value

	case fieldLicense:
		pkg.License = value

	case fieldUrl:
		pkg.Url = value

	default:
		// ignore unknown fields
	}

	return nil
}
