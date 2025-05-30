package preprocess

import (
	"qp/internal/consts"
)

var macroRegistry = map[consts.CmdType][]MacroExpander{
	consts.BlockSelect: {expandSelectMacro},
	consts.BlockWhere:  {expandWhereMacro},
	consts.BlockOrder:  {},
	consts.BlockLimit:  {expandLimitMacro},
}
