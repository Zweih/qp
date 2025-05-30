package driver

import "qp/internal/pkgdata"

type Driver interface {
	Name() string
	Detect() bool
	Load() ([]*pkgdata.PkgInfo, error)
	ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error)
	LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error)
	UpdateHistory(cacheRoot string, pkgs []*pkgdata.PkgInfo) error
	SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error
	IsCacheStale(cacheMtime int64) (bool, error)
}
