package syntax

import (
	"qp/internal/consts"
	"strings"
)

type MacroExpander func(token string) ([]string, bool)

var macroRegistry = map[CmdType][]MacroExpander{
	BlockSelect: {expandSelectMacro},
	BlockWhere:  {expandWhereMacro},
	BlockOrder:  {},
	BlockLimit:  {expandLimitMacro},
}

func macroExpansion(token string, cmd CmdType) ([]string, bool) {
	expanders := macroRegistry[cmd]
	for _, expander := range expanders {
		if replacement, exists := expander(strings.ToLower(token)); exists {
			return replacement, true
		}
	}

	return []string{token}, false
}

func expandSelectMacro(token string) ([]string, bool) {
	switch token {
	case "default":
		return fieldTypesToNames(consts.DefaultFields), true
	case "all":
		return fieldTypesToNames(consts.ValidFields), true
	default:
		return nil, false
	}
}

func fieldTypesToNames(fields []consts.FieldType) []string {
	fieldNames := make([]string, 0, len(fields))
	for _, field := range fields {
		if fieldName, ok := consts.FieldNameLookup[field]; ok {
			fieldNames = append(fieldNames, fieldName)
		}
	}

	return fieldNames
}

func expandWhereMacro(token string) ([]string, bool) {
	switch token {
	case "orphan":
		return []string{"not:required-by", "reason=dependency"}, true
	default:
		return nil, false
	}
}

func expandLimitMacro(token string) ([]string, bool) {
	switch token {
	case "all":
		return []string{"0"}, true
	default:
		return nil, false
	}
}
