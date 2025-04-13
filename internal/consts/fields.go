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
	FieldValidation
	FieldPackager
	FieldDescription
	FieldUrl
	FieldGroups
	FieldSize
	FieldDate
	FieldBuildDate
	FieldVersion
	FieldDepends
	FieldOptDepends
	FieldRequiredBy
	FieldOptionalFor
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
	validation  = "validation"
	packager    = "packager"
	groups      = "groups"
	conflicts   = "conflicts"
	replaces    = "replaces"
	depends     = "depends"
	optdepends  = "optdepends"
	requiredBy  = "required-by"
	optionalFor = "optional-for"
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
	size:        FieldSize,
	pkgType:     FieldPkgType,
	name:        FieldName,
	reason:      FieldReason,
	version:     FieldVersion,
	arch:        FieldArch,
	license:     FieldLicense,
	pkgBase:     FieldPkgBase,
	description: FieldDescription,
	url:         FieldUrl,
	validation:  FieldValidation,
	packager:    FieldPackager,
	groups:      FieldGroups,
	conflicts:   FieldConflicts,
	replaces:    FieldReplaces,
	depends:     FieldDepends,
	optdepends:  FieldOptDepends,
	requiredBy:  FieldRequiredBy,
	optionalFor: FieldOptionalFor,
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
	FieldSize:        size,
	FieldPkgType:     pkgType,
	FieldName:        name,
	FieldReason:      reason,
	FieldVersion:     version,
	FieldArch:        arch,
	FieldLicense:     license,
	FieldPkgBase:     pkgBase,
	FieldDescription: description,
	FieldUrl:         url,
	FieldValidation:  validation,
	FieldPackager:    packager,
	FieldGroups:      groups,
	FieldConflicts:   conflicts,
	FieldReplaces:    replaces,
	FieldDepends:     depends,
	FieldOptDepends:  optdepends,
	FieldRequiredBy:  requiredBy,
	FieldOptionalFor: optionalFor,
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
	// note: this is also the order the columns will be displayed in table output
	ValidFields = []FieldType{
		FieldDate,
		FieldBuildDate,
		FieldSize,
		FieldPkgType,
		FieldName,
		FieldReason,
		FieldVersion,
		FieldArch,
		FieldLicense,
		FieldPkgBase,
		FieldDescription,
		FieldUrl,
		FieldValidation,
		FieldPackager,
		FieldGroups,
		FieldConflicts,
		FieldReplaces,
		FieldDepends,
		FieldOptDepends,
		FieldRequiredBy,
		FieldOptionalFor,
		FieldProvides,
	}
)
