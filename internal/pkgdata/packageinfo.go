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

	Version      string
	Name         string
	ProviderName string
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

	Depends    []Relation
	RequiredBy []Relation
	Provides   []Relation
	Conflicts  []Relation
	Replaces   []Relation
}
