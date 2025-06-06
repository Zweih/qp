package npm

import (
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/storage"
)

type NpmDriver struct {
	modulesDir string
}

func (d *NpmDriver) Name() string {
	return consts.OriginNpm
}

func (d *NpmDriver) Detect() bool {
	modulesDir, err := getGlobalModulesDir()
	if err != nil {
		return false
	}

	if _, err := os.Stat(modulesDir); err != nil {
		return false
	}

	d.modulesDir = modulesDir
	return true
}

func (d *NpmDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name(), d.modulesDir)
}

func (d *NpmDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *NpmDriver) LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(cacheRoot)
}

func (d *NpmDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *NpmDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return shared.BfsStale(d.modulesDir, cacheMtime, 2)
}
