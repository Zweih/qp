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

type Relation struct {
	Depth    int32
	Operator RelationOp

	IsComplex bool

	Name         string
	Version      string
	ProviderName string
	Why          string
	PkgType      string
}

func (rel *Relation) Key() string {
	if rel.PkgType != "" {
		return rel.PkgType + ":" + rel.Name
	}

	return rel.Name
}

func (rel Relation) ProviderKey() string {
	if rel.ProviderName == "" {
		return ""
	}

	if rel.PkgType != "" {
		return rel.PkgType + ":" + rel.ProviderName
	}

	return rel.ProviderName
}

type PkgInfo struct {
	UpdateTimestamp  int64
	BuildTimestamp   int64
	InstallTimestamp int64
	Size             int64
	Name             string
	Reason           string
	Version          string
	Origin           string
	Arch             string
	License          string
	Description      string
	Url              string
	Validation       string
	PkgType          string
	PkgBase          string
	Packager         string

	Groups []string

	Depends     []Relation
	OptDepends  []Relation
	RequiredBy  []Relation
	OptionalFor []Relation
	Provides    []Relation
	Conflicts   []Relation
	Replaces    []Relation
}

func (pkg *PkgInfo) Key() string {
	if pkg.Origin == consts.OriginBrew {
		return pkg.PkgType + ":" + pkg.Name
	}

	return pkg.Name
}

func (pkg *PkgInfo) GetInt(field consts.FieldType) int64 {
	switch field {
	case consts.FieldInstalled:
		return pkg.InstallTimestamp
	case consts.FieldUpdated:
		return pkg.UpdateTimestamp
	case consts.FieldBuilt:
		return pkg.BuildTimestamp
	case consts.FieldSize:
		return pkg.Size
	default:
		panic("invalid field passed to GetInt")
	}
}

func (pkg *PkgInfo) GetString(field consts.FieldType) string {
	switch field {
	case consts.FieldName:
		return pkg.Name
	case consts.FieldReason:
		return pkg.Reason
	case consts.FieldVersion:
		return pkg.Version
	case consts.FieldOrigin:
		return pkg.Origin
	case consts.FieldArch:
		return pkg.Arch
	case consts.FieldLicense:
		return pkg.License
	case consts.FieldDescription:
		return pkg.Description
	case consts.FieldUrl:
		return pkg.Url
	case consts.FieldValidation:
		return pkg.Validation
	case consts.FieldPkgType:
		return pkg.PkgType
	case consts.FieldPkgBase:
		return pkg.PkgBase
	case consts.FieldPackager:
		return pkg.Packager
	default:
		panic("invalid field passed to GetString: " + consts.FieldNameLookup[field])
	}
}

func (pkg *PkgInfo) GetStrArr(field consts.FieldType) []string {
	switch field {
	case consts.FieldGroups:
		return pkg.Groups
	default:
		panic("invalid field passed to GetStringSlice: " + consts.FieldNameLookup[field])
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

func StringToOperator(operatorInput string) RelationOp {
	switch operatorInput {
	case "=", "==":
		return OpEqual
	case "<", "<<":
		return OpLess
	case "<=", "=<":
		return OpLessEqual
	case ">", ">>":
		return OpGreater
	case ">=", "=>":
		return OpGreaterEqual
	default:
		return OpNone
	}
}

func OperatorToString(op RelationOp) string {
	switch op {
	case OpEqual:
		return "="
	case OpLess:
		return "<"
	case OpGreater:
		return ">"
	case OpLessEqual:
		return "<="
	case OpGreaterEqual:
		return ">="
	default:
		return ""
	}
}
