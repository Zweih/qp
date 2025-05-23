package opkg

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type OpkgDriver struct{}

func (d *OpkgDriver) Name() string {
	return consts.OriginOpkg
}

func (d *OpkgDriver) Detect() bool {
	_, err := os.Stat(opkgStatusPath)
	return err == nil
}

func (d *OpkgDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *OpkgDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *OpkgDriver) LoadCache(path string) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path)
}

func (d *OpkgDriver) SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
	return pkgdata.SaveProtoCache(cacheRoot, pkgs)
}

func (d *OpkgDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat(opkgStatusPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix(), nil
}

func (d *OpkgDriver) IsCacheStale(cacheModTime int64) (bool, error) {
	dirInfo, err := os.Stat(opkgStatusPath)
	if err != nil {
		return false, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix() > cacheModTime, nil
}
