package preprocess

import (
	"qp/internal/quipple"
	"strings"
)

func ExpandShortSyntax(args []string) []string {
	expanded := make([]string, 0, len(args))

	for _, arg := range args {
		if cmd, ok := quipple.ShorthandMap[strings.ToLower(arg)]; ok {
			expanded = append(expanded, cmd)
			continue
		}

		expanded = append(expanded, arg)
	}

	return expanded
}
