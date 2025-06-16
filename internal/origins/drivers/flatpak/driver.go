package flatpak

import (
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"qp/internal/storage"
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
			d.installDirs = append(d.installDirs, userPath)
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
	return pkgdata.ResolveDependencyGraph(mergeExtensions(pkgs), nil)
}

func (d *FlatpakDriver) LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
	return storage.LoadProtoCache(cacheRoot)
}

func (d *FlatpakDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return storage.SaveProtoCache(cacheRoot, pkgs)
}

func (d *FlatpakDriver) IsCacheStale(cacheMtime int64) (bool, error) {
	for _, installDir := range d.installDirs {
		appDir := filepath.Join(installDir, typeApp)
		runtimeDir := filepath.Join(installDir, typeRuntime)

		for _, dir := range []string{appDir, runtimeDir} {
			if _, err := os.Stat(dir); err != nil {
				continue
			}

			isStale, err := shared.BfsStale(dir, cacheMtime, 2)
			if err != nil {
				return false, err
			}

			if isStale {
				return true, nil
			}
		}
	}

	return false, nil
}
