package completion

import (
	"qp/internal/quipple"
	"strings"
)

type CmdPair struct {
	Full        string
	Short       string
	Completions string
	Description string
}

func getOrderedCmds() []CmdPair {
	cmdToShorthand := reverseShorthandMap(quipple.ShorthandMap)

	cmdDefs := []struct {
		cmd  string
		desc string
	}{
		{quipple.CmdSelect, "Select fields to display"},
		{quipple.CmdWhere, "Filter packages"},
		{quipple.CmdOrder, "Sort results"},
		{quipple.CmdLimit, "Limit results"},
		{quipple.CmdFormat, "Output format"},
	}

	pairs := make([]CmdPair, 0, len(cmdDefs))
	for _, def := range cmdDefs {
		pairs = append(pairs, CmdPair{
			Full:        def.cmd,
			Short:       cmdToShorthand[def.cmd],
			Description: def.desc,
		})
	}

	return pairs
}

func getAllCmdsStr() string {
	pairs := getOrderedCmds()
	var cmds []string

	for _, pair := range pairs {
		cmds = append(cmds, pair.Full, pair.Short)
	}

	return strings.Join(cmds, " ")
}

func reverseShorthandMap(shorthandMap map[string]string) map[string]string {
	result := make(map[string]string, len(shorthandMap)*2)

	for key, value := range shorthandMap {
		result[value] = key
	}

	return result
}
