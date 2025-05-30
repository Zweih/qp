package pipx

import (
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
)

type PipxDriver struct {
	venvRoot string
}

func (d *PipxDriver) Name() string {
	return consts.OriginPipx
}

func (d *PipxDriver) Detect() bool {
	venvRoot, err := findVenvRoot()
	if err != nil {
		return false
	}

	d.venvRoot = venvRoot
	return true
}

func (d *PipxDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.venvRoot, d.Name())
}

func (d *PipxDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgs, nil
}

func (d *PipxDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *PipxDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *PipxDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return shared.BfsStale(d.venvRoot, cacheMtime, 2)
}
