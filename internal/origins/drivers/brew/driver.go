package brew

import (
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
)

type BrewDriver struct {
	prefix string
}

func (d *BrewDriver) Name() string {
	return consts.OriginBrew
}

func (d *BrewDriver) Detect() bool {
	prefix, err := getPrefix()
	d.prefix = prefix

	return err == nil
}

func (d *BrewDriver) Load(_ string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name(), d.prefix)
}

func (d *BrewDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *BrewDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *BrewDriver) UpdateHistory(_ string, _ []*pkgdata.PkgInfo) error {
	return nil // we don't need to parse logs for brew
}

func (d *BrewDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *BrewDriver) IsCacheStale(cacheModTime int64) (bool, error) {
	cellarPath := filepath.Join(d.prefix, cellarSubPath)
	isCellarStale, err := shared.BfsStale(cellarPath, cacheModTime, 1)
	if err != nil {
		return false, err
	}

	if isCellarStale {
		return true, nil
	}

	caskroomPath := filepath.Join(d.prefix, caskroomSubpath)
	isCaskroomStale, err := shared.BfsStale(caskroomPath, cacheModTime, 1)
	if err != nil {
		return false, err
	}

	return isCaskroomStale, nil
}
