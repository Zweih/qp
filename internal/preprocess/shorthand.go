package preprocess

import (
	"qp/internal/consts"
	"strings"
)

var ShorthandMap = map[string]string{
	"s": consts.CmdSelect,
	"w": consts.CmdWhere,
	"o": consts.CmdOrder,
	"l": consts.CmdLimit,
}

func ExpandShortSyntax(args []string) []string {
	expanded := make([]string, 0, len(args))

	for _, arg := range args {
		if cmd, ok := ShorthandMap[strings.ToLower(arg)]; ok {
			expanded = append(expanded, cmd)
			continue
		}

		expanded = append(expanded, arg)
	}

	return expanded
}
