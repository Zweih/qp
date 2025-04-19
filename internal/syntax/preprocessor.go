package syntax

import "fmt"

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

		expanded, _ := macroExpansion(token, currentBlock)
		if len(expanded) == 0 {
			return nil, fmt.Errorf("macro expansion for %q in block %q produced no output", token, cmdTypeName(currentBlock))
		}

		processed = append(processed, expanded...)
	}

	return processed, nil
}

func cmdTypeName(cmd CmdType) string {
	switch cmd {
	case BlockSelect:
		return CmdSelect
	case BlockWhere:
		return CmdWhere
	case BlockOrder:
		return CmdOrder
	}

	return "[INVALID BLOCK COMMAND]"
}
