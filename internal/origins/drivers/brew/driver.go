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

func (d *BrewDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name(), d.prefix)
}

func (d *BrewDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *BrewDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *BrewDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *BrewDriver) SourceModified() (int64, error) {
	cellarTime, err := shared.GetModTime(filepath.Join(d.prefix, cellarSubPath))
	if err != nil {
		return 0, err
	}

	caskroomTime, err := shared.GetModTime(filepath.Join(d.prefix, caskroomSubpath))
	if err != nil {
		return 0, err
	}

	return max(cellarTime, caskroomTime), nil
}

// TODO: add early return
func (d *BrewDriver) IsCacheStale(cacheModTime int64) (bool, error) {
	cellarTime, err := shared.GetModTime(filepath.Join(d.prefix, cellarSubPath))
	if err != nil {
		return false, err
	}

	caskroomTime, err := shared.GetModTime(filepath.Join(d.prefix, caskroomSubpath))
	if err != nil {
		return false, err
	}

	return max(cellarTime, caskroomTime) > cacheModTime, nil
}
