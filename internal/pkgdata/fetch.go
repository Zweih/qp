package pkgdata

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
	fieldDescription = "%DESC%"
	fieldPkgBase     = "%BASE%"
	fieldReplaces    = "%REPLACES%"
	fieldBuildDate   = "%BUILDDATE%"

	PacmanDbPath = "/var/lib/pacman/local"
)

func FetchPackages() ([]*PkgInfo, error) {
	pkgPaths, err := os.ReadDir(PacmanDbPath)
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
			descPath := filepath.Join(PacmanDbPath, packagePath.Name(), "desc")
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

	// the average desc file is 103.13 lines, reading the entire file into memory is more efficient than using bufio.Scanner
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var pkg PkgInfo
	var currentField string
	start := 0
	end := 0
	length := len(data)

	for end <= length {
		if end == length || data[end] == '\n' {
			line := string(bytes.TrimSpace(data[start:end]))

			switch line {
			case fieldName, fieldInstallDate, fieldSize, fieldReason, fieldVersion,
				fieldArch, fieldLicense, fieldUrl, fieldDescription, fieldPkgBase, fieldBuildDate:
				currentField = line

			case fieldDepends, fieldProvides, fieldConflicts, fieldReplaces:
				currentField = line
				block, next := collectBlockBytes(data, end+1)

				applyMultiLineField(&pkg, currentField, block)
				end = next
				start = next

				continue

			case "":
				currentField = ""

			default:
				if err := applySingleLineField(&pkg, currentField, line); err != nil {
					return nil, fmt.Errorf("error reading desc file %s: %w", descPath, err)
				}
			}

			start = end + 1
		}

		end++
	}

	if pkg.Name == "" {
		return nil, fmt.Errorf("package name is missing in file: %s", descPath)
	}

	if pkg.Reason == "" {
		pkg.Reason = "explicit"
	}

	return &pkg, nil
}

func collectBlockBytes(data []byte, start int) ([]string, int) {
	var block []string
	i := start

	for i < len(data) {
		j := i

		for j < len(data) && data[j] != '\n' {
			j++
		}

		line := bytes.TrimSpace(data[i:j])

		if len(line) == 0 {
			break
		}

		block = append(block, string(line))
		i = j + 1
	}

	return block, i
}

func applySingleLineField(pkg *PkgInfo, field string, value string) error {
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

		pkg.InstallTimestamp = installDate

	case fieldBuildDate:
		buildDate, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid build date value %q: %w", value, err)
		}

		pkg.BuildTimestamp = buildDate

	case fieldVersion:
		pkg.Version = value

	case fieldSize:
		size, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid size value %q: %w", value, err)
		}

		pkg.Size = size

	case fieldArch:
		pkg.Arch = value

	case fieldLicense:
		pkg.License = value

	case fieldUrl:
		pkg.Url = value

	case fieldDescription:
		pkg.Description = value

	case fieldPkgBase:
		pkg.PkgBase = value

	default:
		// ignore unknown fields
	}

	return nil
}

func applyMultiLineField(pkg *PkgInfo, field string, lines []string) {
	relations := parseRelations(lines)

	switch field {
	case fieldDepends:
		pkg.Depends = relations
	case fieldProvides:
		pkg.Provides = relations
	case fieldConflicts:
		pkg.Conflicts = relations
	case fieldReplaces:
		pkg.Replaces = relations
	}
}

func parseRelations(block []string) []Relation {
	relations := make([]Relation, 0, len(block))

	for _, line := range block {
		relations = append(relations, parseRelation(line))
	}

	return relations
}

func parseRelation(input string) Relation {
	opStart := 0
	var depth int32 = 1

	for i := range input {
		switch input[i] {
		case '=', '<', '>':
			opStart = i
			goto parseOp
		}
	}

	return Relation{Name: input, Depth: depth}

parseOp:
	name := input[:opStart]
	opEnd := opStart + 1

	if opEnd < len(input) {
		switch input[opEnd] {
		case '=', '<', '>':
			opEnd++
		}
	}

	operator := stringToOperator(input[opStart:opEnd])
	var version string

	if opEnd < len(input) {
		version = input[opEnd:]
	}

	return Relation{
		Name:     name,
		Operator: operator,
		Version:  version,
		Depth:    depth,
	}
}

func stringToOperator(operatorInput string) RelationOp {
	switch operatorInput {
	case "=":
		return OpEqual
	case "<":
		return OpLess
	case "<=":
		return OpLessEqual
	case ">":
		return OpGreater
	case ">=":
		return OpGreaterEqual
	default:
		return OpNone
	}
}
