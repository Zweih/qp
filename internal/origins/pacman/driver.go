package pacman

import (
	"fmt"
	"os"
	"qp/internal/pkgdata"
)

type PacmanDriver struct{}

func (d *PacmanDriver) Name() string {
	return "pacman"
}

func (d *PacmanDriver) Detect() bool {
	_, err := os.Stat("/var/lib/pacman/local")
	return err == nil
}

func (d *PacmanDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages()
}

func (d *PacmanDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return resolveDependencyGraph(pkgs, nil)
}

func (d *PacmanDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *PacmanDriver) SaveCache(
	path string,
	pkgs []*pkgdata.PkgInfo,
	modTime int64,
) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *PacmanDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat("/var/lib/pacman/local")
	if err != nil {
		return 0, fmt.Errorf("failed to read pacman DB mod time: %v", err)
	}

	return dirInfo.ModTime().Unix(), nil
}
