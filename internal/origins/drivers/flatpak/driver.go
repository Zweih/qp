package flatpak

import (
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type FlatpakDriver struct {
	installDirs []string
}

func (d *FlatpakDriver) Name() string {
	return consts.OriginFlatpak
}

func (d *FlatpakDriver) Detect() bool {
	if _, err := os.Stat(systemInstallDir); err == nil {
		d.installDirs = append(d.installDirs, systemInstallDir)
	}

	home, err := os.UserHomeDir()
	if err == nil {
		userPath := filepath.Join(home, userInstallDir)
		if _, err := os.Stat(userPath); err == nil {
			d.installDirs = append(d.installDirs, userInstallDir)
		}
	}

	if len(d.installDirs) == 0 {
		return false
	}

	return true
}

func (d *FlatpakDriver) Load(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name(), d.installDirs)
}

func (d *FlatpakDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *FlatpakDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return []*pkgdata.PkgInfo{}, nil
}

func (d *FlatpakDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return nil
}

func (d *FlatpakDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	return true, nil
}
