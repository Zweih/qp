package pkgdata

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	fieldName        = "%NAME%"
	fieldInstallDate = "%INSTALLDATE%"
	fieldSize        = "%SIZE%"
	fieldReason      = "%REASON%"
)

func FetchPackages() ([]PackageInfo, error) {
	pacmanDbPath := "/var/lib/pacman/local"
	// entries instead of dirs since there can be files or directories
	entries, err := os.ReadDir(pacmanDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pacman database: %v", err)
	}

	var packages []PackageInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		descPath := filepath.Join(pacmanDbPath, entry.Name(), "desc")
		pkg, err := parseDescFile(descPath)

		if err == nil {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
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
		case fieldName, fieldInstallDate, fieldSize, fieldReason:
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
