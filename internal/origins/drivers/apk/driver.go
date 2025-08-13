package apk

import (
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/storage"
)

type ApkDriver struct{}

func (d *ApkDriver) Name() string {
	return consts.OriginApk
}

func (d *ApkDriver) Detect() bool {
	_, err := os.Stat(apkDbPath)
	return err == nil
}

func (d *ApkDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *ApkDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *ApkDriver) LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(cacheRoot)
}

func (d *ApkDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *ApkDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	isDbStale, err := shared.IsDirStale(d.Name(), apkDbPath, cacheMtime)
	if err != nil {
		return false, err
	}

	if isDbStale {
		return true, nil
	}

	isWorldStale, err := shared.IsDirStale(d.Name(), apkWorldPath, cacheMtime)
	if err != nil {
		return false, nil
	}

	return isWorldStale, nil
}
