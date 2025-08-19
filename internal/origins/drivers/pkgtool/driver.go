package pkgtool

import (
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/storage"
)

type PkgtoolDriver struct{}

func (d *PkgtoolDriver) Name() string {
	return consts.OriginPkgtool
}

func (d *PkgtoolDriver) Detect() bool {
	_, err := os.Stat(packagesDbPath)

	return err == nil
}

func (d *PkgtoolDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *PkgtoolDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *PkgtoolDriver) LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(cacheRoot)
}

func (d *PkgtoolDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *PkgtoolDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return shared.BfsStale(packagesDbPath, cacheMtime, 1)
}
