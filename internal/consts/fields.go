package consts

type FieldType int

// ordered by filter efficiency
const (
	FieldReason FieldType = iota
	FieldArch
	FieldName
	FieldOrigin
	FieldEnv
	FieldPkgType
	FieldLicense
	FieldPkgBase
	FieldValidation
	FieldPackager
	FieldDescription
	FieldUrl
	FieldAlsoIn
	FieldGroups
	FieldSize
	FieldInstalled
	FieldUpdated
	FieldBuilt
	FieldVersion
	FieldDepends
	FieldOptDepends
	FieldRequiredBy
	FieldOptionalFor
	FieldProvides
	FieldConflicts
	FieldReplaces
)

type FieldPrim int32

const (
	FieldPrimDate = iota
	FieldPrimSize
	FieldPrimStr
	FieldPrimStrArr
	FieldPrimRel
)

func GetFieldPrim(field FieldType) FieldPrim {
	switch field {
	case FieldUpdated, FieldBuilt, FieldInstalled:
		return FieldPrimDate

	case FieldSize:
		return FieldPrimSize

	case FieldName, FieldReason, FieldVersion,
		FieldOrigin, FieldArch, FieldLicense,
		FieldUrl, FieldDescription, FieldValidation,
		FieldEnv, FieldPkgType, FieldPkgBase, FieldPackager:
		return FieldPrimStr

	case FieldAlsoIn, FieldGroups:
		return FieldPrimStrArr

	case FieldConflicts, FieldReplaces, FieldDepends,
		FieldOptDepends, FieldRequiredBy, FieldOptionalFor,
		FieldProvides:
		return FieldPrimRel

	default:
		panic("invalid field passed to GetFieldPrim")
	}
}

const (
	installed   = "installed"
	updated     = "updated"
	built       = "built"
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
	env         = "env"
	pkgType     = "pkgtype"
	pkgBase     = "pkgbase"
	packager    = "packager"
	alsoIn      = "also-in"
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
	"u":    FieldUpdated,
	"n":    FieldName,
	"r":    FieldReason,
	"s":    FieldSize,
	"v":    FieldVersion,
	"D":    FieldDepends,
	"R":    FieldRequiredBy,
	"p":    FieldProvides,
	"bd":   FieldBuilt,
	"type": FieldPkgType,

	"date":         FieldUpdated, // legacy field
	"build-date":   FieldBuilt,   // legacy field
	"alphabetical": FieldName,    // legacy flag, to be deprecated

	installed:   FieldInstalled,
	updated:     FieldUpdated,
	built:       FieldBuilt,
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
	env:         FieldEnv,
	pkgType:     FieldPkgType,
	pkgBase:     FieldPkgBase,
	packager:    FieldPackager,
	alsoIn:      FieldAlsoIn,
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
	FieldInstalled:   installed,
	FieldUpdated:     updated,
	FieldBuilt:       built,
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
	FieldEnv:         env,
	FieldPkgType:     pkgType,
	FieldPkgBase:     pkgBase,
	FieldPackager:    packager,
	FieldAlsoIn:      alsoIn,
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
		FieldUpdated,
		FieldName,
		FieldReason,
		FieldSize,
	}
	// note: this is also the order the columns will be displayed in table output
	ValidFields = []FieldType{
		FieldUpdated,
		FieldBuilt,
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
		FieldEnv,
		FieldPkgType,
		FieldPkgBase,
		FieldPackager,
		FieldAlsoIn,
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
