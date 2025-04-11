package consts

type (
	FieldType    int
	SubfieldType int32
)

// ordered by filter efficiency
const (
	FieldReason FieldType = iota
	FieldPkgType
	FieldArch
	FieldLicense
	FieldName
	FieldPkgBase
	FieldDescription
	FieldUrl
	FieldSize
	FieldDate
	FieldBuildDate
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
	buildDate   = "build-date"
	name        = "name"
	reason      = "reason"
	size        = "size"
	version     = "version"
	pkgType     = "pkgtype"
	arch        = "arch"
	license     = "license"
	pkgBase     = "pkgbase"
	description = "description"
	url         = "url"
	conflicts   = "conflicts"
	replaces    = "replaces"
	depends     = "depends"
	requiredBy  = "required-by"
	provides    = "provides"

	target = "target"
	depth  = "depth"
)

var FieldTypeLookup = map[string]FieldType{
	"d":    FieldDate,
	"n":    FieldName,
	"r":    FieldReason,
	"s":    FieldSize,
	"v":    FieldVersion,
	"D":    FieldDepends,
	"R":    FieldRequiredBy,
	"p":    FieldProvides,
	"bd":   FieldBuildDate,
	"type": FieldPkgType,

	"alphabetical": FieldName, // legacy flag, to be deprecated

	date:        FieldDate,
	buildDate:   FieldBuildDate,
	name:        FieldName,
	reason:      FieldReason,
	size:        FieldSize,
	version:     FieldVersion,
	pkgType:     FieldPkgType,
	arch:        FieldArch,
	license:     FieldLicense,
	pkgBase:     FieldPkgBase,
	description: FieldDescription,
	url:         FieldUrl,
	conflicts:   FieldConflicts,
	replaces:    FieldReplaces,
	depends:     FieldDepends,
	requiredBy:  FieldRequiredBy,
	provides:    FieldProvides,
}

var SubfieldTypeLookup = map[string]SubfieldType{
	"":     SubfieldTarget,
	target: SubfieldTarget,
	depth:  SubfieldDepth,
}

var FieldNameLookup = map[FieldType]string{
	FieldDate:        date,
	FieldBuildDate:   buildDate,
	FieldName:        name,
	FieldSize:        size,
	FieldReason:      reason,
	FieldVersion:     version,
	FieldPkgType:     pkgType,
	FieldArch:        arch,
	FieldLicense:     license,
	FieldPkgBase:     pkgBase,
	FieldDescription: description,
	FieldUrl:         url,
	FieldConflicts:   conflicts,
	FieldReplaces:    replaces,
	FieldDepends:     depends,
	FieldRequiredBy:  requiredBy,
	FieldProvides:    provides,
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
		FieldBuildDate,
		FieldName,
		FieldReason,
		FieldSize,
		FieldVersion,
		FieldPkgType,
		FieldArch,
		FieldLicense,
		FieldPkgBase,
		FieldDescription,
		FieldUrl,
		FieldConflicts,
		FieldReplaces,
		FieldDepends,
		FieldRequiredBy,
		FieldProvides,
	}
)
