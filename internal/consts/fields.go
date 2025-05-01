package consts

type FieldType int

// ordered by filter efficiency
const (
	FieldReason FieldType = iota
	FieldArch
	FieldName
	FieldOrigin
	FieldPkgType
	FieldLicense
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
	date        = "date"
	buildDate   = "build-date"
	size        = "size"
	name        = "name"
	reason      = "reason"
	version     = "version"
	origin      = "origin"
	arch        = "arch"
	license     = "license"
	description = "description"
	url         = "url"
	validation  = "validation"
	pkgType     = "pkgtype"
	pkgBase     = "pkgbase"
	packager    = "packager"
	groups      = "groups"
	conflicts   = "conflicts"
	replaces    = "replaces"
	depends     = "depends"
	optdepends  = "optdepends"
	requiredBy  = "required-by"
	optionalFor = "optional-for"
	provides    = "provides"
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
	name:        FieldName,
	reason:      FieldReason,
	version:     FieldVersion,
	origin:      FieldOrigin,
	arch:        FieldArch,
	license:     FieldLicense,
	description: FieldDescription,
	url:         FieldUrl,
	validation:  FieldValidation,
	pkgType:     FieldPkgType,
	pkgBase:     FieldPkgBase,
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

var FieldNameLookup = map[FieldType]string{
	FieldDate:        date,
	FieldBuildDate:   buildDate,
	FieldSize:        size,
	FieldName:        name,
	FieldReason:      reason,
	FieldVersion:     version,
	FieldOrigin:      origin,
	FieldArch:        arch,
	FieldLicense:     license,
	FieldDescription: description,
	FieldUrl:         url,
	FieldValidation:  validation,
	FieldPkgType:     pkgType,
	FieldPkgBase:     pkgBase,
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
		FieldName,
		FieldReason,
		FieldVersion,
		FieldOrigin,
		FieldArch,
		FieldLicense,
		FieldDescription,
		FieldUrl,
		FieldValidation,
		FieldPkgType,
		FieldPkgBase,
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

var StringFields = map[FieldType]struct{}{
	FieldName:        {},
	FieldReason:      {},
	FieldVersion:     {},
	FieldArch:        {},
	FieldLicense:     {},
	FieldDescription: {},
	FieldUrl:         {},
	FieldValidation:  {},
	FieldPkgType:     {},
	FieldPkgBase:     {},
	FieldPackager:    {},
}

var RelationFields = map[FieldType]struct{}{
	FieldDepends:     {},
	FieldOptDepends:  {},
	FieldRequiredBy:  {},
	FieldOptionalFor: {},
	FieldProvides:    {},
	FieldConflicts:   {},
	FieldReplaces:    {},
}

var RangeFields = map[FieldType]struct{}{
	FieldDate:      {},
	FieldBuildDate: {},
	FieldSize:      {},
}
