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

func GetCompletions() string {
	var fieldNames []string
	for fieldName := range consts.FieldTypeLookup {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	completionData := prepareCompletionData(fieldNames)
	cmds := getOrderedCmds()

	var bashCases []string
	for _, cmd := range cmds {
		data := completionData[cmd.Full]
		bashCase := generateBashCase(cmd, data.bashCompletions)

		if cmd.Full == quipple.CmdWhere {
			bashCase += `
      if [[ "${cur}" != *"=" ]] && [[ "${COMPREPLY[0]}" == *"=" ]]; then
        compopt -o nospace 2>/dev/null || true
      fi
      return 0
      ;;`
		} else {
			bashCase += `
      return 0
      ;;`
		}

		bashCases = append(bashCases, bashCase)
	}

	var zshCases []string
	for _, cmd := range cmds {
		data := completionData[cmd.Full]
		zshCase := generateZshCase(cmd, data.zshCompletions, data.description)

		if cmd.Full == quipple.CmdSelect {
			zshCase += ` -S ','`
		} else if cmd.Full == quipple.CmdWhere {
			zshCase += ` -S ''`
		}

		zshCase += `
        ;;`

		zshCases = append(zshCases, zshCase)
	}

	var zshArgCases []string
	for _, cmd := range cmds {
		data := completionData[cmd.Full]
		zshCase := generateZshCase(cmd, data.zshCompletions, data.description)

		if cmd.Full == quipple.CmdSelect {
			zshCase += ` -S ','`
		} else if cmd.Full == quipple.CmdWhere {
			zshCase += ` -S ''`
		}

		zshCase += `
          ;;`

		zshArgCases = append(zshArgCases, zshCase)
	}

	bashScript := fmt.Sprintf(`
if [[ -n "$BASH_VERSION" ]]; then
  _qp_completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD - 1]}"
    local prev_lower="${prev,,}"

    case "${prev_lower}" in
%s
    *)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      return 0
      ;;
    esac
  }
  complete -F _qp_completion qp
fi`,
		strings.Join(bashCases, "\n"),
		getAllCmdsStr(),
	)

	zshScript := fmt.Sprintf(`
if [[ -n "$ZSH_VERSION" ]]; then
  _qp_completion() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    case $words[1] in
    qp)
      if (( CURRENT == 2 )); then
%s
      else
        local cmd_lower="${words[2]:l}"
        case $cmd_lower in
%s
        esac
      fi
      ;;
    esac
  }
  compdef _qp_completion qp
fi`,
		generateZshCmdValues(),
		strings.Join(zshArgCases, "\n"),
	)

	return fmt.Sprintf("%s\n%s", bashScript, zshScript)
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
