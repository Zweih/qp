package preprocess

import (
	"fmt"
	"qp/internal/quipple"
	"strings"
)

func Preprocess(args []string) ([]string, error) {
	tokens := normalizeWhitespace(args)
	tokens = ExpandShortSyntax(tokens)

	var processed []string
	currentBlock := quipple.BlockNone

	for _, token := range tokens {
		if block := quipple.CmdTypeLookup[strings.ToLower(token)]; block != quipple.BlockNone {
			currentBlock = block
			processed = append(processed, token)
			continue
		}

		if currentBlock == quipple.BlockNone {
			return nil, fmt.Errorf("unexpected token: %q (expected a command block like 'select', 'where', or 'order')", token)
		}

		if currentBlock == quipple.BlockWhere {
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
				return nil, fmt.Errorf("macro expansion for %q in block %q produced no output", subtoken, quipple.CmdNameLookup[currentBlock])
			}

			processed = append(processed, expanded...)
		}
	}

	return processed, nil
}

func normalizeWhitespace(args []string) []string {
	var tokens []string
	for _, arg := range args {
		rawTokens := strings.FieldsFunc(arg, func(r rune) bool {
			return r == ' ' || r == '\n'
		})

		for _, token := range rawTokens {
			if len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
	}

	return tokens
}

func getSubtokens(block quipple.CmdType, token string) []string {
	if block == quipple.BlockSelect && strings.Contains(token, ",") {
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
