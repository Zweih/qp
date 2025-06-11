package deb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	"sort"
	"strings"
)

var licenseHints = []struct {
	phrase  string
	license string
	weight  int
}{
	{"most of the gnu c library", "LGPL-2.1+", 10},
	{"gnu c library is free software", "LGPL-2.1+", 10},
	{"gnu lesser general public license", "LGPL-2.1+", 8},
	{"gnu general public license", "GPL-2+", 6},
	{"distributed under the bsd license", "BSD", 5},
	{"bsd license", "BSD", 5},
	{"bsd-4-clause", "BSD-4-Clause", 4},
	{"mit license", "MIT", 3},
	{"expat", "Expat", 3},
	{"isc license", "ISC", 3},
	{"public domain", "public-domain", 2},
	{"lgpl", "LGPL", 2},
}

func extractLicense(pkg *pkgdata.PkgInfo) error {
	licensePath, err := resolveLicensePath(pkg.Name)
	if err != nil {
		pkg.License = unknownLicense
		return err
	}

	file, err := os.Open(licensePath)
	if err != nil {
		pkg.License = unknownLicense
		return fmt.Errorf("failed to open license file for %s: %w", pkg.Name, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		pkg.License = unknownLicense
		return fmt.Errorf("failed to load license file for %s: %w", pkg.Name, err)
	}

	var fallbackLicense string

	for block := range bytes.SplitSeq(data, []byte("\n\n")) {
		var isFilesAll bool
		var license string

		for byteLine := range bytes.SplitSeq(block, []byte("\n")) {
			line := string(byteLine)
			if strings.HasPrefix(line, filesPrefix) {
				isFilesAll = strings.TrimSpace(strings.TrimPrefix(line, filesPrefix)) == "*"
			}

			if strings.HasPrefix(line, licensePrefix) {
				license = strings.TrimSpace(strings.TrimPrefix(line, licensePrefix))
				if fallbackLicense == "" {
					fallbackLicense = license
				}
			}
		}

		if isFilesAll && license != "" {
			pkg.License = license
			return nil
		}
	}

	if match, ok := matchLicenseText(data); ok {
		pkg.License = match
		return nil
	}

	if fallbackLicense != "" {
		pkg.License = fallbackLicense
		return nil
	}

	pkg.License = "custom"
	return nil
}

func resolveLicensePath(packageName string) (string, error) {
	basePath := filepath.Join(licensePath, packageName, licenseFileName)
	resolvedPath, err := filepath.EvalSymlinks(basePath)

	if err == nil && fileExists(resolvedPath) {
		return resolvedPath, nil
	}

	symlinkTarget, statErr := os.Readlink(basePath)
	var fallbackPrefix string

	if statErr == nil {
		parts := strings.Split(symlinkTarget, string(filepath.Separator))
		if len(parts) >= 2 {
			fallbackPrefix = parts[len(parts)-2]
		}
	}

	if fallbackPrefix == "" {
		fallbackPrefix = packageName
	}

	pattern := filepath.Join(licensePath, fallbackPrefix+"*/"+licenseFileName)
	matches, _ := filepath.Glob(pattern)

	if len(matches) == 0 {
		return "", fmt.Errorf("no license file found for package %s", packageName)
	}

	sort.Strings(matches)
	return matches[0], nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func matchLicenseText(data []byte) (string, bool) {
	text := strings.ToLower(string(data))

	bestMatch := ""
	bestWeight := 0

	for _, hint := range licenseHints {
		if strings.Contains(text, hint.phrase) && hint.weight > bestWeight {
			bestMatch = hint.license
			bestWeight = hint.weight
		}
	}

	return bestMatch, bestMatch != ""
}
