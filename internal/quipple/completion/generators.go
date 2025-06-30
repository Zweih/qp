package completion

import (
	"fmt"
	"strings"
)

func generateBashCase(pair CmdPair, completions string) string {
	return fmt.Sprintf(`    %s | %s)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))`,
		pair.Full, pair.Short, completions)
}

func generateZshCase(pair CmdPair, completions string, description string) string {
	tag := strings.ReplaceAll(description, " ", "-")

	return fmt.Sprintf(`      %s | %s)
        local -a opts
        opts=(%s)
        _describe -t %s '%s' opts`,
		pair.Full, pair.Short, completions, tag, description)
}

func generateZshCmdValues() string {
	cmds := getOrderedCmds()
	var values []string

	values = append(values, "_values 'commands' \\")

	for i, cmd := range cmds {
		fullCmd := fmt.Sprintf("          '%s[%s]'", cmd.Full, cmd.Description)
		values = append(values, fullCmd)

		shortDesc := fmt.Sprintf("%s (short)", cmd.Description)
		shortCmd := fmt.Sprintf("          '%s[%s]'", cmd.Short, shortDesc)

		if i < len(cmds)-1 {
			shortCmd += " \\"
		}

		values = append(values, shortCmd)

		values[len(values)-2] += " \\"
	}

	return strings.Join(values, "\n")
}
