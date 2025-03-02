package pkgdata

import "time"

type BasePackageInfo struct {
	Name       string   `json:"name,omitempty"`
	Reason     string   `json:"reason,omitempty"`  // "explicit" or "dependency"
	Size       int64    `json:"size,omitempty"`    // package size in bytes
	Version    string   `json:"version,omitempty"` // current installed version
	Depends    []string `json:"depends,omitempty"`
	RequiredBy []string `json:"requiredBy,omitempty"`
	Provides   []string `json:"provides,omitempty"`
}

// info about a single installed package
type PackageInfo struct {
	Timestamp time.Time
	BasePackageInfo
}

type PackageInfoJson struct {
	Timestamp *time.Time `json:"timestamp,omitempty"`
	BasePackageInfo
}
