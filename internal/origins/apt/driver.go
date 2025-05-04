package apt

import (
	"fmt"
	"os"
	"qp/internal/pkgdata"
)

type AptDriver struct{}

func (d *AptDriver) Name() string {
	return "apt"
}

func (d *AptDriver) Detect() bool {
	_, err := os.Stat(dpkgPath)
	return err == nil
}

func (d *AptDriver) Load() ([]*pkgdata.PkgInfo, error) {
	return fetchPackages(d.Name())
}

func (d *AptDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.ResolveDependencyGraph(pkgs, nil)
}

func (d *AptDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *AptDriver) SaveCache(
	path string,
	pkgs []*pkgdata.PkgInfo,
	modTime int64,
) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *AptDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat(dpkgPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix(), nil
}
