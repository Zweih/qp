package quipple

type CmdType int

const (
	BlockNone CmdType = iota
	BlockSelect
	BlockWhere
	BlockOrder
	BlockLimit
	BlockFormat
)

const (
	CmdSelect = "select"
	CmdWhere  = "where"
	CmdOrder  = "order"
	CmdLimit  = "limit"
	CmdFormat = "format"
)

var CmdTypeLookup = map[string]CmdType{
	CmdSelect: BlockSelect,
	CmdWhere:  BlockWhere,
	CmdOrder:  BlockOrder,
	CmdLimit:  BlockLimit,
	CmdFormat: BlockFormat,
}

var CmdNameLookup = map[CmdType]string{
	BlockSelect: CmdSelect,
	BlockWhere:  CmdWhere,
	BlockOrder:  CmdOrder,
	BlockLimit:  CmdLimit,
	BlockFormat: CmdFormat,
}

const (
	MacroOrphan      = "orphan"
	MacroSuperOrphan = "superorphan"
	MacroHeavy       = "heavy"
	MacroLight       = "light"
	MacroAll         = "all"
	MacroDefault     = "default"
)

var SelectMacros = []string{
	MacroAll,
	MacroDefault,
}

var WhereMacros = []string{
	MacroOrphan,
	MacroSuperOrphan,
	MacroHeavy,
	MacroLight,
}

var LimitMacros = []string{
	MacroAll,
}
