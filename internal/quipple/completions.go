package quipple

import (
	"fmt"
	"qp/internal/consts"
	"sort"
	"strings"
)

func formatForBash(completions []string) string {
	return strings.Join(completions, " ")
}

func formatForZsh(completions []string) string {
	return "'" + strings.Join(completions, "' '") + "'"
}

func GetBashCompletion() string {
	commands := []string{
		CmdSelect, "s",
		CmdWhere, "w",
		CmdOrder, "o",
		CmdLimit, "l",
		CmdFormat, "f",
	}
	commandsStr := strings.Join(commands, " ")

	var fieldNames []string
	for fieldName := range consts.FieldTypeLookup {
		fieldNames = append(fieldNames, fieldName)
	}

	selectCompletions := append(fieldNames, SelectMacros...)
	sort.Strings(selectCompletions)

	var whereFieldNames []string
	for _, fieldName := range fieldNames {
		whereFieldNames = append(whereFieldNames, fieldName+"=")
	}
	sort.Strings(whereFieldNames)

	var whereCompletions []string
	whereCompletions = append(whereCompletions, whereFieldNames...)
	whereCompletions = append(whereCompletions, WhereMacros...)
	sort.Strings(whereCompletions)

	orderCompletions := fieldNames
	limitCompletions := LimitMacros
	formatCompletions := []string{consts.OutputTable, consts.OutputJSON, consts.OutputKeyValue}

	selectBash := formatForBash(selectCompletions)
	whereBash := formatForBash(whereCompletions)
	orderBash := formatForBash(orderCompletions)
	limitBash := formatForBash(limitCompletions)
	formatBash := formatForBash(formatCompletions)

	selectZsh := formatForZsh(selectCompletions)
	whereZsh := formatForZsh(whereCompletions)
	orderZsh := formatForZsh(orderCompletions)
	limitZsh := formatForZsh(limitCompletions)
	formatZsh := formatForZsh(formatCompletions)

	return fmt.Sprintf(
		`
if [[ -n "$BASH_VERSION" ]]; then
  _qp_completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD - 1]}"

    case "${prev}" in
    %s | s)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      return 0
      ;;
    %s | w)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      if [[ "${cur}" != *"=" ]] && [[ "${COMPREPLY[0]}" == *"=" ]]; then
        compopt -o nospace 2>/dev/null || true
      fi
      return 0
      ;;
    %s | o)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      return 0
      ;;
    %s | l)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      return 0
      ;;
    %s | f)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      return 0
      ;;
    *)
      COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      return 0
      ;;
    esac
  }
  complete -F _qp_completion qp
fi

if [[ -n "$ZSH_VERSION" ]]; then
  _qp_completion() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    case $words[1] in
    qp)
      case $words[2] in
      select | s)
        local -a fields
        fields=(%s)
        _describe -t fields 'fields' fields -S ','
        ;;
      where | w)
        local -a where_opts
        where_opts=(%s)
        _describe -t where-options 'where options' where_opts -S ''
        ;;
      order | o)
        local -a order_opts
        order_opts=(%s)
        _describe -t order-options 'order options' order_opts
        ;;
      limit | o)
        local -a limit_opts
        limit_opts=(%s)
        _describe -t limit-options 'limit options' limit_opts
        ;;
      format | f)
        local -a formats
        formats=(%s)
        _describe -t formats 'formats' formats
        ;;
      *)
        _values 'commands' \
          'select[Select fields to display]' \
          's[Select fields (short)]' \
          'where[Filter packages]' \
          'w[Where (short)]' \
          'order[Sort results]' \
          'o[Order (short)]' \
          'limit[Limit results]' \
          'l[Limit (short)]' \
          'format[Output format]' \
          'f[Format (short)]'
        ;;
      esac
      ;;
    esac
  }
  compdef _qp_completion qp
fi
`,
		CmdSelect,
		selectBash,
		CmdWhere,
		whereBash,
		CmdOrder,
		orderBash,
		CmdLimit,
		limitBash,
		CmdFormat,
		formatBash,
		commandsStr,
		selectZsh,
		whereZsh,
		orderZsh,
		limitZsh,
		formatZsh,
	)
}
