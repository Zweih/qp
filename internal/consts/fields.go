package consts

type FieldType int

const (
	FieldDate FieldType = iota
	FieldName
	FieldReason
	FieldSize
	FieldVersion
	FieldDepends
	FieldRequiredBy
	FieldProvides
	FieldConflicts
	FieldArch
	FieldLicense
	FieldUrl
)

const (
	date       = "date"
	name       = "name"
	reason     = "reason"
	size       = "size"
	version    = "version"
	depends    = "depends"
	requiredBy = "required-by"
	provides   = "provides"
	conflicts  = "conflicts"
	arch       = "arch"
	license    = "license"
	url        = "url"
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

	date:       FieldDate,
	name:       FieldName,
	reason:     FieldReason,
	size:       FieldSize,
	version:    FieldVersion,
	depends:    FieldDepends,
	requiredBy: FieldRequiredBy,
	provides:   FieldProvides,
	conflicts:  FieldConflicts,
	arch:       FieldArch,
	license:    FieldLicense,
	url:        FieldUrl,
}

var FieldNameLookup = map[FieldType]string{
	FieldDate:       date,
	FieldName:       name,
	FieldReason:     reason,
	FieldVersion:    version,
	FieldDepends:    depends,
	FieldRequiredBy: requiredBy,
	FieldProvides:   provides,
	FieldConflicts:  conflicts,
	FieldArch:       arch,
	FieldLicense:    license,
	FieldUrl:        url,
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
	}
)
