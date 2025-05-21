package brew

import (
	"qp/internal/consts"
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

func (d *BrewDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *BrewDriver) SaveCache(
	path string,
	pkgs []*pkgdata.PkgInfo,
	modTime int64,
) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *BrewDriver) SourceModified() (int64, error) {
	cellarTime, err := getModTime(d.prefix, cellarSubPath)
	if err != nil {
		return 0, err
	}

	caskroomTime, err := getModTime(d.prefix, caskroomSubpath)
	if err != nil {
		return 0, err
	}

	return max(cellarTime, caskroomTime), nil
}
