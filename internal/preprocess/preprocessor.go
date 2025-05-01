package preprocess

import (
	"fmt"
	"qp/internal/consts"
	"strings"
)

func Preprocess(args []string) ([]string, error) {
	args = ExpandShortSyntax(args)

	var processed []string
	currentBlock := consts.BlockNone

	for _, token := range args {
		if block := consts.CmdTypeLookup[token]; block != consts.BlockNone {
			currentBlock = block
			processed = append(processed, token)
			continue
		}

		if currentBlock == consts.BlockNone {
			return nil, fmt.Errorf("unexpected token: %q (expected in a command block like 'select', 'where', or 'order')", token)
		}

		if currentBlock == consts.BlockWhere {
			normalized := normalizeNegationShorthand(token)
			if len(normalized) == 2 {
				processed = append(processed, normalized...)
				continue
			}
		}

		subtokens := getSubtokens(currentBlock, token)
		for _, subtoken := range subtokens {
			subtoken = strings.TrimSpace(subtoken)
			expanded, _ := macroExpansion(subtoken, currentBlock)

			if len(expanded) == 0 {
				return nil, fmt.Errorf("macro expansion for %q in block %q produced no output", subtoken, consts.CmdNameLookup[currentBlock])
			}

			processed = append(processed, expanded...)
		}
	}

	return processed, nil
}

func getSubtokens(block consts.CmdType, token string) []string {
	if block == consts.BlockSelect && strings.Contains(token, ",") {
		return strings.Split(token, ",")
	}

	return []string{token}
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
