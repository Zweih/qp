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
	for _, iPkg := range installed {
		wanted[iPkg.Name] = struct{}{}
	}

	metadata, err := loadMetadata(caskCachePath, getCaskKey, wanted)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	result := make([]*pkgdata.PkgInfo, 0, len(installed))
	for _, iPkg := range installed {
		pkg := &pkgdata.PkgInfo{
			Name:    iPkg.Name,
			Version: iPkg.Version,
			PkgType: typeCask,
			Origin:  origin,
		}

		if meta, ok := metadata[pkg.Name]; ok {
			mergeCaskMetadata(pkg, meta)
		}

		result = append(result, pkg)
	}

	return result, nil
}

func getInstalledCasks(caskroomRoot string) ([]installedPkg, error) {
	entries, err := os.ReadDir(caskroomRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read Caskroom directory: %w", err)
	}

	var pkgs []installedPkg
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		subEntries, err := os.ReadDir(filepath.Join(caskroomRoot, name))
		if err != nil {
			return nil, fmt.Errorf("failed to read cask %s directory: %w", name, err)
		}

		var version string
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() || subEntry.Name() == ".metadata" {
				continue
			}
			version := subEntry.Name()
			fmt.Printf("Cask: %s, Version: %s\n", name, version)
			break
		}

		pkgs = append(pkgs, installedPkg{
			Name:    name,
			Version: version,
		})
	}

	return pkgs, nil
}
