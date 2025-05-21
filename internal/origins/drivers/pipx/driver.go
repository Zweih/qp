package pipx

import (
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"time"
)

type PipxDriver struct {
	venvRooot string
}

func (d *PipxDriver) Name() string {
	return consts.OriginPipx
}

func (d *PipxDriver) Detect() bool {
	d.venvRooot = getVenvRoot()
	_, err := os.Stat(d.venvRooot)

	return err == nil
}

func getVenvRoot() string {
	if custom := os.Getenv("PIPX_HOME"); custom != "" {
		return filepath.Join(custom, "venvs")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}

	return filepath.Join(home, defaultVenvRoot)
}

func (d *PipxDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.venvRooot, d.Name())
}

func (d *PipxDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgs, nil
}

func (d *PipxDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *PipxDriver) SaveCache(path string, pkgs []*pkgdata.PkgInfo, modTime int64) error {
	return nil
}

func (d *PipxDriver) SourceModified() (int64, error) {
	return time.Now().Unix() + 200, nil
}
