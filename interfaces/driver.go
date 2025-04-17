package interfaces

import "qp/internal/pkgdata"

type Driver interface {
	Name() string
	Detect() bool
	Load() ([]*pkgdata.PkgInfo, error)
	ResolveDeps(pkgs []*pkgdata.PkgInfo) ([]*pkgdata.PkgInfo, error)
	LoadCache(path string, modTime int64) ([]*pkgdata.PkgInfo, error)
	SaveCache(path string, pkgs []*pkgdata.PkgInfo, modTime int64) error
	SourceModified() (int64, error)
}
