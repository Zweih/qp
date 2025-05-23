package pipx

import (
	"errors"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"time"
)

type PipxDriver struct {
	venvRoot string
}

func (d *PipxDriver) Name() string {
	return consts.OriginPipx
}

func (d *PipxDriver) Detect() bool {
	venvRoot, err := getVenvRoot()
	if err != nil {
		return false
	}

	d.venvRoot = venvRoot
	return true
}

func getVenvRoot() (string, error) {
	var venvRootPath string

	if custom := os.Getenv("PIPX_HOME"); custom != "" {
		venvRootPath = filepath.Join(custom, "venvs")
		_, err := os.Stat(venvRootPath)
		if err == nil {
			return venvRootPath, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}

	possibleRoots := []string{defaultVenvRoot, otherVenvRoot}
	for _, root := range possibleRoots {
		venvRootPath = filepath.Join(home, root)
		_, err := os.Stat(venvRootPath)
		if err == nil {
			return venvRootPath, nil
		}
	}

	return "", errors.New("no pipx venv root found")
}

func (d *PipxDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.venvRoot, d.Name())
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
