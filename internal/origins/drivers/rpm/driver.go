package rpm

import (
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
)

type RpmDriver struct {
	dbPath string
}

func (d *RpmDriver) Name() string {
	return consts.OriginRpm
}

func (d *RpmDriver) Detect() bool {
	possibleRoots := []string{
		defaultRpmRoot,
		modernRpmRoot,
		rebuildRpmRoot,
	}

	for _, pRoot := range possibleRoots {
		if fullPath := detectRpmDatabase(pRoot); fullPath != "" {
			d.dbPath = fullPath
			return true
		}
	}

	return false
}

func (d *RpmDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name(), filepath.Join(d.dbPath))
}

func (d *RpmDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, evaluateComplexDependency)
}

func (d *RpmDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *RpmDriver) UpdateHistory(_ string, _ []*pkgdata.PkgInfo) error {
	return nil
}

func (d *RpmDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *RpmDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return shared.IsDirStale(d.Name(), d.dbPath, cacheMtime)
}
