package preprocess

import (
	"qp/internal/consts"
	"qp/internal/quipple"
	"strings"
)

type MacroExpander func(token string) ([]string, bool)

func macroExpansion(token string, cmd quipple.CmdType) ([]string, bool) {
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
	case quipple.MacroDefault:
		return fieldTypesToNames(consts.DefaultFields), true
	case quipple.MacroAll:
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
	var expanded []string

	switch token {
	case quipple.MacroOrphan:
		expanded = []string{"no:required-by", "and", "reason=dependency"}
	case quipple.MacroSuperOrphan:
		expanded = []string{"no:required-by", "and", "reason=dependency", "and", "no:optional-for"}
	case quipple.MacroHeavy:
		expanded = []string{"size=100MB:"}
	case quipple.MacroLight:
		expanded = []string{"size=:1MB"}
	default:
		return nil, false
	}

	return append([]string{"q"}, append(expanded, "p")...), true
}

func expandLimitMacro(token string) ([]string, bool) {
	var prefix string
	if strings.ContainsRune(token, ':') {
		parts := strings.SplitAfter(token, ":")
		prefix, token = parts[0], parts[1]
	}

	switch token {
	case quipple.MacroAll:
		return []string{prefix + "0"}, true
	default:
		return nil, false
	}
}
