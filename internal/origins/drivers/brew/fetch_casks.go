package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
)

func fetchCasks(
	origin string,
	prefix string,
) ([]*pkgdata.PkgInfo, error) {
	caskroomRoot := filepath.Join(prefix, caskroomSubpath)
	installed, err := getInstalledCasks(caskroomRoot)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	wanted := make(map[string]struct{}, len(installed))
	for _, name := range installed {
		wanted[name] = struct{}{}
	}

	metadata, err := loadMetadata(caskCachePath, getCaskKey, wanted)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	result := make([]*pkgdata.PkgInfo, 0, len(installed))
	for _, name := range installed {
		receiptPath := filepath.Join(caskroomRoot, name, ".metadata", receiptName)
		pkg, err := parseCaskReceipt(receiptPath)
		if err != nil {
			continue
		}

		pkg.Name = name
		pkg.Origin = origin

		size, err := getInstallSize(filepath.Join(caskroomRoot, name, pkg.Version))
		if err == nil {
			pkg.Size = size
		}

		fmt.Println(pkg.Depends)

		if meta, ok := metadata[pkg.Name]; ok {
			mergeCaskMetadata(pkg, meta)
		}

		result = append(result, pkg)
	}

	return result, nil
}

func getInstalledCasks(caskroomRoot string) ([]string, error) {
	entries, err := os.ReadDir(caskroomRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read Caskroom directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		names = append(names, entry.Name())
	}

	return names, nil
}
