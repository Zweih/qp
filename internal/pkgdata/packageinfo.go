package pkgdata

type RelationOp int32

const (
	OpNone RelationOp = iota
	OpEqual
	OpLess
	OpLessEqual
	OpGreater
	OpGreaterEqual
)

type PkgType int32

const (
	PkgTypeUnknown PkgType = iota
	PkgTypePkg
	PkgTypeSplit
	PkgTypeSrc
	PkgTypeDebug
)

type Relation struct {
	Depth    int32
	Operator RelationOp

	Name         string
	Version      string
	ProviderName string
	Why          string
}

type PkgInfo struct {
	InstallTimestamp int64
	BuildTimestamp   int64
	Size             int64

	PkgType PkgType

	Name        string
	Reason      string
	Version     string
	Arch        string
	License     string
	Url         string
	Description string
	PkgBase     string
	Validation  string
	Packager    string

	Groups []string

	Depends    []Relation
	OptDepends []Relation
	RequiredBy []Relation
	Provides   []Relation
	Conflicts  []Relation
	Replaces   []Relation
}
