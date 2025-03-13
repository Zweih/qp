package consts

type FieldType string

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
)

const (
	FieldDate       FieldType = date
	FieldName       FieldType = name
	FieldReason     FieldType = reason
	FieldSize       FieldType = size
	FieldVersion    FieldType = version
	FieldDepends    FieldType = depends
	FieldRequiredBy FieldType = requiredBy
	FieldProvides   FieldType = provides
	FieldConflicts  FieldType = conflicts
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
	}
)
