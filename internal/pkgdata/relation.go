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

type Relation struct {
	Depth    int32
	Operator RelationOp

	IsComplex bool

	Name         string
	Version      string
	ProviderName string
	Why          string
	PkgType      string // pkgtype should only be declared when there is a split ecosystem, such as formulae/casks in Brew
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
