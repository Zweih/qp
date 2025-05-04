package apt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
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
	path := filepath.Join(licensePath, pkg.Name, licenseFileName)
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to read license file for %s: %w", pkg.Name, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
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

	return fmt.Errorf("no license found for %s", pkg.Name)
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
