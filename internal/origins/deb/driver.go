package deb

import (
	"fmt"
	"os"
	"qp/internal/pkgdata"
)

type DebDriver struct{}

func (d *DebDriver) Name() string {
	return "deb"
}

func (d *DebDriver) Detect() bool {
	_, err := os.Stat(dpkgPath)
	return err == nil
}

func (d *DebDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *DebDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *DebDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *DebDriver) SaveCache(
	path string,
	pkgs []*pkgdata.PkgInfo,
	modTime int64,
) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *DebDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat(dpkgPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix(), nil
}
