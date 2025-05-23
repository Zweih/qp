package pacman

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type PacmanDriver struct{}

func (d *PacmanDriver) Name() string {
	return consts.OriginPacman
}

func (d *PacmanDriver) Detect() bool {
	_, err := os.Stat(PacmanDbPath)
	return err == nil
}

func (d *PacmanDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *PacmanDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *PacmanDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *PacmanDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *PacmanDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat(PacmanDbPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read pacman DB mod time: %v", err)
	}

	return dirInfo.ModTime().Unix(), nil
}

func (d *PacmanDriver) IsCacheStale(cacheModTime int64) (bool, error) {
	dirInfo, err := os.Stat(PacmanDbPath)
	if err != nil {
		return false, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix() > cacheModTime, nil
}
