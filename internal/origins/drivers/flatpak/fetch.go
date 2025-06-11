package flatpak

import (
	"bufio"
	"fmt"
	"os"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"
	"sync"
)

func fetchPackages(origin string, installDirs []string) ([]*pkgdata.PkgInfo, error) {
	pkgRefs, err := discoverPackages(installDirs)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	if len(pkgRefs) < 1 {
		return []*pkgdata.PkgInfo{}, nil
	}

	inputChan := make(chan *PkgRef, len(pkgRefs))
	errChan := make(chan error, worker.DefaultBufferSize)
	var errGroup sync.WaitGroup

	for _, pkgRef := range pkgRefs {
		inputChan <- pkgRef
	}
	close(inputChan)

	return []*pkgdata.PkgInfo{}, nil
}

func parseMetadata(pkgRef *PkgRef) (*pkgdata.PkgInfo, error) {
	file, err := os.Open(pkgRef.MetadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s metadta file %s: %w", pkgRef.Name, pkgRef.MetadataPath, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix("#", line) {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}

		if strings.Contains(line, "=") {
			key, value := parseKeyValue(line)
		}
	}
}

func parseKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func applyMetadataField(
	pkg *pkgdata.PkgInfo,
	section string,
	key string,
	value string,
) error {
	switch section {
	case sectionApplication:
		return applyApplicationField
	default:
		// TODO: extension
		return nil
	}
}

func applyApplicationField(pkg *pkgdata.PkgInfo, key string, value string) error {
  switch key

}
