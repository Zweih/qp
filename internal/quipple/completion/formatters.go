package completion

import "strings"

func formatForBash(completions []string) string {
	return strings.Join(completions, " ")
}

func formatForZsh(completions []string) string {
	return "'" + strings.Join(completions, "' '") + "'"
}
