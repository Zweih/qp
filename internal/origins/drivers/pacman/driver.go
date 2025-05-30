package pacman

import (
	"os"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
)

type PacmanDriver struct{}

func (d *PacmanDriver) Name() string {
	return consts.OriginPacman
}

func (d *PacmanDriver) Detect() bool {
	_, err := os.Stat(pacmanDbDir)
	return err == nil
}

func (d *PacmanDriver) Load(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name(), cacheRoot)
}

func (d *PacmanDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *PacmanDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *PacmanDriver) UpdateHistory(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return nil
}

func (d *PacmanDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *PacmanDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return shared.IsDirStale(d.Name(), pacmanDbDir, cacheMtime)
}
