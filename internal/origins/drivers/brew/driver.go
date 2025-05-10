package brew

import (
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
	return fetchPackages(d.Name(), d.prefix)
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
	info, err := os.Stat(filepath.Join(d.prefix, cellarSubPath))
	if err != nil {
		return 0, err
	}

	return info.ModTime().Unix(), nil
}
