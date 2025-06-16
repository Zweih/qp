package driver

import "qp/internal/pkgdata"

type Driver interface {
	Name() string
	Detect() bool
	Load(cachePath string) ([]*pkgdata.PkgInfo, error)
	ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error)
	LoadCache(cachePath string) ([]*pkgdata.PkgInfo, error)
	SaveCache(cachePath string, pkgs []*pkgdata.PkgInfo) error
	IsCacheStale(cacheMtime int64) (bool, error)
}
