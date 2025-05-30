package opkg

import (
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/storage"
)

type OpkgDriver struct{}

func (d *OpkgDriver) Name() string {
	return consts.OriginOpkg
}

func (d *OpkgDriver) Detect() bool {
	_, err := os.Stat(opkgStatusPath)
	return err == nil
}

func (d *OpkgDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *OpkgDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *OpkgDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(path)
}

func (d *OpkgDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *OpkgDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return shared.IsDirStale(d.Name(), opkgStatusPath, cacheMtime)
}
