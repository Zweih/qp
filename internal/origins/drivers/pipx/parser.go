package pipx

import (
	"bufio"
	"os"
	"qp/internal/pkgdata"
	"strings"
)

// TODO: parse line by line
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
		case "Name":
			pkg.Name = value
		case "Version":
			pkg.Version = value
		case "Summary":
			pkg.Description = value
		case "License":
			pkg.License = value
		case "Home-page":
			if homepage == "" {
				homepage = value
			}
		case "Project-URL":
			fieldParts := strings.SplitN(value, ",", 2)
			if len(fieldParts) != 2 {
				continue
			}

			if strings.TrimSpace(fieldParts[0]) == "Homepage" {
				homepage = strings.TrimSpace(fieldParts[1])
			}
		}
	}

	pkg.Url = homepage

	return pkg, nil
}
