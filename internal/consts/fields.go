package consts

type (
	FieldType    int
	SubfieldType int32
)

type FilterKey struct {
	Field    FieldType
	Subfield SubfieldType
}

// ordered by filter efficiency
const (
	FieldReason FieldType = iota
	FieldArch
	FieldLicense
	FieldName
	FieldPkgBase
	FieldDescription
	FieldUrl
	FieldSize
	FieldDate
	FieldVersion
	FieldDepends
	FieldRequiredBy
	FieldProvides
	FieldConflicts
)

const (
	SubfieldNone SubfieldType = iota
	SubfieldName
	SubfieldDepth
)

const (
	date        = "date"
	name        = "name"
	reason      = "reason"
	size        = "size"
	version     = "version"
	description = "description"
	pkgBase     = "pkgbase"
	depends     = "depends"
	requiredBy  = "required-by"
	provides    = "provides"
	conflicts   = "conflicts"
	arch        = "arch"
	license     = "license"
	url         = "url"
	depth       = "depth"
)

var FieldTypeLookup = map[string]FieldType{
	"d": FieldDate,
	"n": FieldName,
	"r": FieldReason,
	"s": FieldSize,
	"v": FieldVersion,
	"D": FieldDepends,
	"R": FieldRequiredBy,
	"p": FieldProvides,

	"alphabetical": FieldName, // legacy flag, to be deprecated

	date:        FieldDate,
	name:        FieldName,
	reason:      FieldReason,
	arch:        FieldArch,
	license:     FieldLicense,
	url:         FieldUrl,
	description: FieldDescription,
	pkgBase:     FieldPkgBase,
	size:        FieldSize,
	version:     FieldVersion,
	depends:     FieldDepends,
	requiredBy:  FieldRequiredBy,
	provides:    FieldProvides,
	conflicts:   FieldConflicts,
}

var SubfieldTypeLookup = map[string]SubfieldType{
	name:  SubfieldName,
	depth: SubfieldDepth,
}

var FieldNameLookup = map[FieldType]string{
	FieldDate:        date,
	FieldName:        name,
	FieldSize:        size,
	FieldReason:      reason,
	FieldVersion:     version,
	FieldDepends:     depends,
	FieldRequiredBy:  requiredBy,
	FieldProvides:    provides,
	FieldConflicts:   conflicts,
	FieldArch:        arch,
	FieldLicense:     license,
	FieldUrl:         url,
	FieldDescription: description,
	FieldPkgBase:     pkgBase,
}

var SubfieldNameLookup = map[SubfieldType]string{
	SubfieldName:  name,
	SubfieldDepth: depth,
}

var (
	DefaultFields = []FieldType{
		FieldDate,
		FieldName,
		FieldReason,
		FieldSize,
	}
	ValidFields = []FieldType{
		FieldDate,
		FieldName,
		FieldReason,
		FieldSize,
		FieldVersion,
		FieldDepends,
		FieldRequiredBy,
		FieldProvides,
		FieldConflicts,
		FieldArch,
		FieldLicense,
		FieldUrl,
		FieldDescription,
		FieldPkgBase,
	}
)
