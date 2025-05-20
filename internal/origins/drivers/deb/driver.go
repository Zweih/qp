package deb

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
)

type DebDriver struct {
	reasonMap map[string]string
}

func (d *DebDriver) Name() string {
	return consts.OriginDeb
}

func (d *DebDriver) Detect() bool {
	_, err := os.Stat(dpkgPath)
	return err == nil
}

func (d *DebDriver) Load() ([]*pkgdata.PkgInfo, error) {
	reasonMap, _ := loadInstallReasons()
	d.reasonMap = reasonMap

	return fetchPackages(d.Name(), reasonMap)
}

func (d *DebDriver) ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error) {
	// won't be needed again, free memory
	defer func() {
		d.reasonMap = nil
	}()

	resolvedPkgs, err := pkgdata.ResolveDependencyGraph(pkgs, nil)
	if err != nil {
		return resolvedPkgs, err
	}

	isFile := false
	var modTime int64

	if info, err := os.Stat(installReasonPath); err == nil {
		isFile = true
		modTime = info.ModTime().Unix()
	}

	for _, pkg := range resolvedPkgs {
		_, hasReason := d.reasonMap[pkg.Name]

		if !hasReason && pkg.Reason == consts.ReasonExplicit && len(pkg.RequiredBy) > 0 {
			if isFile && modTime > pkg.InstallTimestamp {
				pkg.Reason = consts.ReasonDependency
			}
		}
	}

	return resolvedPkgs, nil
}

func (d *DebDriver) LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error) {
	return pkgdata.LoadProtoCache(path, modTime)
}

func (d *DebDriver) SaveCache(path string, pkgs []*pkgdata.PkgInfo, modTime int64) error {
	return pkgdata.SaveProtoCache(pkgs, path, modTime)
}

func (d *DebDriver) SourceModified() (int64, error) {
	dirInfo, err := os.Stat(dpkgPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s DB mod time: %v", d.Name(), err)
	}

	return dirInfo.ModTime().Unix(), nil
}
