package pkgdata

type PkgInfo struct {
	Timestamp  int64    `json:"timestamp,omitempty"`
	Name       string   `json:"name,omitempty"`
	Reason     string   `json:"reason,omitempty"`  // "explicit" or "dependency"
	Size       int64    `json:"size,omitempty"`    // package size in bytes
	Version    string   `json:"version,omitempty"` // current installed version
	Depends    []string `json:"depends,omitempty"`
	RequiredBy []string `json:"requiredBy,omitempty"`
	Provides   []string `json:"provides,omitempty"`
	Conflicts  []string `json:"conflicts,omitempty"`
	Arch       string   `json:"arch,omitempty"`
	License    string   `json:"license,omitempty"`
	Url        string   `json:"url,omitempty"`
}

func convertToPointers(pkgs []PkgInfo) []*PkgInfo {
	pkgPtrs := make([]*PkgInfo, len(pkgs))
	for i := range pkgs {
		pkgPtrs[i] = &pkgs[i]
	}

	return pkgPtrs
}

func dereferencePkgPointers(pkgPtrs []*PkgInfo) []PkgInfo {
	pkgs := make([]PkgInfo, len(pkgPtrs))
	for i := range pkgPtrs {
		pkgs[i] = *pkgPtrs[i]
	}

	return pkgs
}
