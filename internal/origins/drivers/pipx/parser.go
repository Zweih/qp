package pipx

import (
	"bufio"
	"fmt"
	"os"
	"qp/internal/pkgdata"
	"strings"
)

func parseMetadataFile(metadataPath string) (*pkgdata.PkgInfo, error) {
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	pkg := &pkgdata.PkgInfo{}
	scanner := bufio.NewScanner(file)

	var homepage string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case fieldName:
			pkg.Name = value
		case fieldVersion:
			pkg.Version = value
		case fieldSummary:
			pkg.Description = value
		case fieldLicense:
			pkg.License = value
		case fieldHomepage:
			if homepage == "" {
				homepage = value
			}
		case fieldProjectUrl:
			fieldParts := strings.SplitN(value, ",", 2)
			if len(fieldParts) != 2 {
				continue
			}

			if strings.TrimSpace(fieldParts[0]) == subfieldHomepage {
				homepage = strings.TrimSpace(fieldParts[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading metadata file %s: %w", metadataPath, err)
	}

	pkg.Url = homepage

	return pkg, nil
}

func parseWheelFile(wheelPath string) (string, error) {
	file, err := os.Open(wheelPath)
	if err != nil {
		return "", err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "Tag" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("no tag field in wheel file for %s", wheelPath)
}
