package driver

import "qp/internal/pkgdata"

type Driver interface {
	Name() string
	Detect() bool
	Load() ([]*pkgdata.PkgInfo, error)
	ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error)
	LoadCache(cacheRoot string) ([]*pkgdata.PkgInfo, error)
	SaveCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error
	SourceModified() (int64, error)
	IsCacheStale(cacheModTime int64) (bool, error)
}
