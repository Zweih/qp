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
	Name     string
	Version  string
	Operator RelationOp
}

type PkgInfo struct {
	Timestamp   int64
	Size        int64
	Name        string
	Reason      string
	Version     string
	Arch        string
	License     string
	Url         string
	Description string
	PkgBase     string
	Depends     []Relation
	RequiredBy  []Relation
	Provides    []Relation
	Conflicts   []Relation
}
