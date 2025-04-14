package pkgdata

import (
	"qp/internal/consts"
)

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
	PkgBase     string
	Description string
	Url         string
	Validation  string
	Packager    string

	Groups []string

	Depends     []Relation
	OptDepends  []Relation
	RequiredBy  []Relation
	OptionalFor []Relation
	Provides    []Relation
	Conflicts   []Relation
	Replaces    []Relation
}

func (pkg *PkgInfo) GetString(field consts.FieldType) string {
	switch field {
	case consts.FieldName:
		return pkg.Name
	case consts.FieldReason:
		return pkg.Reason
	case consts.FieldVersion:
		return pkg.Version
	case consts.FieldArch:
		return pkg.Arch
	case consts.FieldLicense:
		return pkg.License
	case consts.FieldPkgBase:
		return pkg.PkgBase
	case consts.FieldDescription:
		return pkg.Description
	case consts.FieldUrl:
		return pkg.Url
	case consts.FieldValidation:
		return pkg.Validation
	case consts.FieldPackager:
		return pkg.Packager
	default:
		panic("invalid field passed to GetString: " + consts.FieldNameLookup[field])
	}
}

func (pkg *PkgInfo) GetRelations(field consts.FieldType) []Relation {
	switch field {
	case consts.FieldConflicts:
		return pkg.Conflicts
	case consts.FieldReplaces:
		return pkg.Replaces
	case consts.FieldDepends:
		return pkg.Depends
	case consts.FieldOptDepends:
		return pkg.OptDepends
	case consts.FieldRequiredBy:
		return pkg.RequiredBy
	case consts.FieldOptionalFor:
		return pkg.OptionalFor
	case consts.FieldProvides:
		return pkg.Provides
	default:
		panic("invalid field passed to GetRelations")
	}
}

func (pkg *PkgInfo) GetInt(field consts.FieldType) int64 {
	switch field {
	case consts.FieldDate:
		return pkg.InstallTimestamp
	case consts.FieldBuildDate:
		return pkg.BuildTimestamp
	case consts.FieldSize:
		return pkg.Size
	case consts.FieldPkgType:
		return int64(pkg.PkgType)
	default:
		panic("invalid field passed to GetInt")
	}
}
