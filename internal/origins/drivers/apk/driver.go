package apk

import (
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type ApkDriver struct {
	modulesDirs []string
}

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
	return []*pkgdata.PkgInfo{}, nil
}

func (d *ApkDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return nil
}

func (d *ApkDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return true, nil
}
