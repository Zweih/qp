package opkg

import (
	"fmt"
	"os"
	"qp/internal/pkgdata"
)

type OpkgDriver struct{}

func (d *OpkgDriver) Name() string {
	return "opkg"
}

func (d *OpkgDriver) Detect() bool {
	_, err := os.Stat(opkgStatusPath)
	return err == nil
}

func (d *OpkgDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *OpkgDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *OpkgDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *OpkgDriver) SaveCache(path string, pkgs []*pkgdata.PkgInfo, modTime int64) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *OpkgDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat(opkgStatusPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix(), nil
}
