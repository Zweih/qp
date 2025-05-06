package deb

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type DebDriver struct {
	fallbackNeeded bool
}

func (d *DebDriver) Name() string {
	return "deb"
}

func (d *DebDriver) Detect() bool {
	_, err := os.Stat(dpkgPath)
	return err == nil
}

func (d *DebDriver) Load() ([]*pkgdata.PkgInfo, error) {
	reasonMap, err := loadInstallReasons()
	d.fallbackNeeded = err != nil

	return fetchPackages(d.Name(), reasonMap)
}

func (d *DebDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	resolvedPkgs, err := pkgdata.ResolveDependencyGraph(pkgs, nil)
	if err != nil {
		return resolvedPkgs, err
	}

	if d.fallbackNeeded {
		for _, pkg := range resolvedPkgs {
			if pkg.Reason == consts.ReasonExplicit && len(pkg.RequiredBy) > 0 {
				pkg.Reason = consts.ReasonDependency
			}
		}
	}

	return resolvedPkgs, nil
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
