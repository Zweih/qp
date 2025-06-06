package storage

const (
	cacheVersion   = 24 // bump when updating structure of PkgInfo/Relation/pkginfo.proto OR when dependency resolution is updated
	historyVersion = 3

	xdgCacheHomeEnv = "XDG_CACHE_HOME"
	homeEnv         = "HOME"
	sudoUserEnv     = "SUDO_USER"
	userEnv         = "USER"

	qpCacheDir = "query-packages"

	dotCache   = ".cache"
	dotModTime = ".modtime"
	dotHistory = ".history"
	dotLock    = ".lock"

	darwinCacheDir = "Library/Caches"
)
