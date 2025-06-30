package completion

import (
	"fmt"
	"qp/internal/consts"
	"qp/internal/quipple"
	"sort"
	"strings"
)

type CompletionData struct {
	bashCompletions string
	zshCompletions  string
	description     string
}

func GetCompletions(shell string) (string, error) {
	var fieldNames []string
	for fieldName := range consts.FieldTypeLookup {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	completionData := prepareCompletionData(fieldNames)
	cmds := getOrderedCmds()

	switch shell {
	case "bash":
		return getBashCompletion(cmds, completionData), nil
	case "zsh":
		return getZshCompletion(cmds, completionData), nil
	}

	return "", fmt.Errorf("no completions available for: %s", shell)
}

func getBashCompletion(cmds []CmdPair, completionData map[string]CompletionData) string {
	var bashCases []string
	for _, cmd := range cmds {
		data := completionData[cmd.Full]
		bashCase := generateBashCase(cmd, data.bashCompletions)

		switch cmd.Full {
		case quipple.CmdSelect:
			bashCase += fmt.Sprintf(bashSelectCase, data.bashCompletions, data.bashCompletions)
		case quipple.CmdWhere:
			bashCase += bashWhereCase
		default:
			bashCase += bashDefaultCase
		}

		bashCases = append(bashCases, bashCase)
	}

	return fmt.Sprintf(
		bashScriptTemplate,
		strings.Join(bashCases, "\n"),
		getAllCmdsStr(),
	)
}

func getZshCompletion(cmds []CmdPair, completionData map[string]CompletionData) string {
	var zshCases []string
	for _, cmd := range cmds {
		data := completionData[cmd.Full]
		zshCase := generateZshCase(cmd, data.zshCompletions, data.description)

		switch cmd.Full {
		case quipple.CmdSelect:
			zshCase = fmt.Sprintf(
				zshSelectCase,
				cmd.Full,
				cmd.Short,
				data.zshCompletions,
				data.zshCompletions,
				strings.ReplaceAll(data.description, " ", "-"),
				data.description,
			)
		case quipple.CmdWhere:
			zshCase += zshWhereCase
		}

		zshCase += zshCaseSuffix

		zshCases = append(zshCases, zshCase)
	}

	return fmt.Sprintf(
		zshScriptTemplate,
		generateZshCmdValues(),
		strings.Join(zshCases, "\n"),
	)
}

func createCompletionData(completions []string, description string, shouldSort bool) CompletionData {
	if shouldSort {
		sort.Strings(completions)
	}

	return CompletionData{
		bashCompletions: formatForBash(completions),
		zshCompletions:  formatForZsh(completions),
		description:     description,
	}
}

func prepareCompletionData(fieldNames []string) map[string]CompletionData {
	data := make(map[string]CompletionData)

	selectCompletions := append(fieldNames, quipple.SelectMacros...)

	whereFieldNames := make([]string, len(fieldNames))
	for i, fieldName := range fieldNames {
		whereFieldNames[i] = fieldName + "="
	}
	whereCompletions := append(whereFieldNames, quipple.WhereMacros...)

	formatCompletions := []string{consts.OutputTable, consts.OutputJSON, consts.OutputKeyValue}

	data[quipple.CmdSelect] = createCompletionData(selectCompletions, "fields", true)
	data[quipple.CmdWhere] = createCompletionData(whereCompletions, "where options", true)
	data[quipple.CmdOrder] = createCompletionData(fieldNames, "order options", false)
	data[quipple.CmdLimit] = createCompletionData(quipple.LimitMacros, "limit options", false)
	data[quipple.CmdFormat] = createCompletionData(formatCompletions, "formats", false)

	return data
}
