package pkgdata

type RelationOp int

const (
	OpNone RelationOp = iota
	OpEqual
	OpLess
	OpLessEqual
	OpGreater
	OpGreaterEqual
)

type Relation struct {
	Name     string     `json:"name"`
	Version  string     `json:"version,omitempty"`
	Operator RelationOp `json:"operator,omitempty"`
}

type PkgInfo struct {
	Timestamp  int64      `json:"timestamp,omitempty"`
	Size       int64      `json:"size,omitempty"` // package size in bytes
	Name       string     `json:"name,omitempty"`
	Reason     string     `json:"reason,omitempty"`  // "explicit" or "dependency"
	Version    string     `json:"version,omitempty"` // current installed version
	Arch       string     `json:"arch,omitempty"`
	License    string     `json:"license,omitempty"`
	Url        string     `json:"url,omitempty"`
	Depends    []Relation `json:"depends,omitempty"`
	RequiredBy []Relation `json:"requiredBy,omitempty"`
	Provides   []Relation `json:"provides,omitempty"`
	Conflicts  []Relation `json:"conflicts,omitempty"`
}
