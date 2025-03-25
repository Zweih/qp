package pkgdata

type PkgInfo struct {
	Timestamp  int64    `json:"timestamp,omitempty"`
	Size       int64    `json:"size,omitempty"` // package size in bytes
	Name       string   `json:"name,omitempty"`
	Reason     string   `json:"reason,omitempty"`  // "explicit" or "dependency"
	Version    string   `json:"version,omitempty"` // current installed version
	Arch       string   `json:"arch,omitempty"`
	License    string   `json:"license,omitempty"`
	Url        string   `json:"url,omitempty"`
	Depends    []string `json:"depends,omitempty"`
	RequiredBy []string `json:"requiredBy,omitempty"`
	Provides   []string `json:"provides,omitempty"`
	Conflicts  []string `json:"conflicts,omitempty"`
}
