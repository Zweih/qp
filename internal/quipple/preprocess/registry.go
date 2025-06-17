package preprocess

import (
	"qp/internal/quipple"
)

var macroRegistry = map[quipple.CmdType][]MacroExpander{
	quipple.BlockSelect: {expandSelectMacro},
	quipple.BlockWhere:  {expandWhereMacro},
	quipple.BlockOrder:  {},
	quipple.BlockLimit:  {expandLimitMacro},
	quipple.BlockFormat: {expandFormatMacro},
}
