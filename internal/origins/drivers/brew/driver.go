package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
)

type BrewDriver struct {
	prefix string
}

func (d *BrewDriver) Name() string {
	return "brew"
}

func (d *BrewDriver) Detect() bool {
	prefix, err := getPrefix()
	d.prefix = prefix

	return err == nil
}

func (d *BrewDriver) Load() ([]*pkgdata.PkgInfo, error) {
	formulae, err := fetchFormulae(d.Name(), d.prefix)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	casks, err := fetchCasks(d.Name(), d.prefix)
	if err != nil {
		return []*pkgdata.PkgInfo{}, err
	}

	return append(formulae, casks...), nil
}

func (d *BrewDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *BrewDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *BrewDriver) SaveCache(
	path string,
	pkgs []*pkgdata.PkgInfo,
	modTime int64,
) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *BrewDriver) SourceModified() (int64, error) {
	cellarPath := filepath.Join(d.prefix, cellarSubPath)
	entries, err := os.ReadDir(cellarPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read Cellar: %w", err)
	}

	var latest int64
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := os.Stat(filepath.Join(cellarPath, entry.Name()))
		if err != nil {
			continue
		}

		modTime := info.ModTime().Unix()
		if modTime > latest {
			latest = modTime
		}
	}

	cellarInfo, err := os.Stat(cellarPath)
	if err != nil {
		return 0, fmt.Errorf("failed to stat Cellar: %w", err)
	}

	modTime := cellarInfo.ModTime().Unix()
	if modTime > latest {
		latest = modTime
	}

	return latest, nil
}
