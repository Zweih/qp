package consts

type (
	FieldType    int
	SubfieldType int32
)

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
	FieldReplaces
)

const (
	SubfieldDepth SubfieldType = iota
	SubfieldTarget
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
	replaces    = "replaces"
	arch        = "arch"
	license     = "license"
	url         = "url"
	target      = "target"
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
	replaces:    FieldReplaces,
}

var SubfieldTypeLookup = map[string]SubfieldType{
	"":     SubfieldTarget,
	target: SubfieldTarget,
	depth:  SubfieldDepth,
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
	FieldReplaces:    replaces,
	FieldArch:        arch,
	FieldLicense:     license,
	FieldUrl:         url,
	FieldDescription: description,
	FieldPkgBase:     pkgBase,
}

var SubfieldNameLookup = map[SubfieldType]string{
	SubfieldTarget: target,
	SubfieldDepth:  depth,
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
		FieldPkgBase,
		FieldArch,
		FieldLicense,
		FieldDescription,
		FieldUrl,
		FieldConflicts,
		FieldReplaces,
		FieldDepends,
		FieldRequiredBy,
		FieldProvides,
	}
)
