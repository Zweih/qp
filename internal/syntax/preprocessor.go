package syntax

import (
	"fmt"
	"strings"
)

func Preprocess(args []string) ([]string, error) {
	args = ExpandShortSyntax(args)

	var processed []string
	currentBlock := BlockNone

	for _, token := range args {
		if block := lookupCommand(token); block != BlockNone {
			currentBlock = block
			processed = append(processed, token)
			continue
		}

		if currentBlock == BlockNone {
			return nil, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', or 'order')", token)
		}

		if currentBlock == BlockWhere {
			normalized := normalizeNegationShorthand(token)
			if len(normalized) == 2 {
				processed = append(processed, normalized...)
				continue
			}
		}

		expanded, _ := macroExpansion(token, currentBlock)
		if len(expanded) == 0 {
			return nil, fmt.Errorf("macro expansion for %q in block %q produced no output", token, cmdTypeName(currentBlock))
		}

		processed = append(processed, expanded...)
	}

	return processed, nil
}

func normalizeNegationShorthand(input string) []string {
	if strings.Contains(input, "!=") || strings.Contains(input, "!==") {
		replaced := strings.ReplaceAll(input, "!==", "==")
		replaced = strings.ReplaceAll(replaced, "!=", "==")
		return []string{"not", replaced}
	}

	if strings.HasPrefix(input, "no:") {
		field := strings.TrimPrefix(input, "no:")
		return []string{"not", "has:" + field}
	}

	return []string{input}
}

func cmdTypeName(cmd CmdType) string {
	switch cmd {
	case BlockSelect:
		return CmdSelect
	case BlockWhere:
		return CmdWhere
	case BlockOrder:
		return CmdOrder
	case BlockLimit:
		return CmdLimit
	}

	return "[INVALID BLOCK COMMAND]"
}
