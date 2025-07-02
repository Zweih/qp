package snap

import (
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"qp/internal/storage"
)

type SnapDriver struct{}

func (d *SnapDriver) Name() string {
	return consts.OriginSnap
}

func (d *SnapDriver) Detect() bool {
	entries, err := os.ReadDir(snapRoot)
	if err != nil {
		return false
	}

	hasSnaps := false
	for _, entry := range entries {
		if entry.IsDir() {
			hasSnaps = true
			break
		}
	}

	return hasSnaps
}

func (d *SnapDriver) Load(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages()
}

func (d *SnapDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *SnapDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(path)
}

func (d *SnapDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *SnapDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	// return shared.IsDirStale(d.Name(), pacmanDbDir, cacheMtime)
	return true, nil
}
