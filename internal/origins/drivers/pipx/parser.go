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
			pkg.Url = value
		}
	}

	return pkg, nil
}
