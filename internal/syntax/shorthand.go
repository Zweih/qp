package syntax

import "strings"

var ShorthandMap = map[string]string{
	"s": CmdSelect,
	"w": CmdWhere,
	"o": CmdOrder,
	"l": CmdLimit,
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
