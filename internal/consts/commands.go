package consts

type CmdType int

const (
	BlockNone CmdType = iota
	BlockSelect
	BlockWhere
	BlockOrder
	BlockLimit
)

const (
	CmdSelect = "select"
	CmdWhere  = "where"
	CmdOrder  = "order"
	CmdLimit  = "limit"
)

var CmdTypeLookup = map[string]CmdType{
	CmdSelect: BlockSelect,
	CmdWhere:  BlockWhere,
	CmdOrder:  BlockOrder,
	CmdLimit:  BlockLimit,
}

var CmdNameLookup = map[CmdType]string{
	BlockSelect: CmdSelect,
	BlockWhere:  CmdWhere,
	BlockOrder:  CmdOrder,
	BlockLimit:  CmdLimit,
}
