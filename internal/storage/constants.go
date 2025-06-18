package storage

const (
	// bump when updating structure of PkgInfo/Relation/pkginfo.proto OR when you've changed something that requires the cache to be rebuilt upon install
	cacheVersion   = 32
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
