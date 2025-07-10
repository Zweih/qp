package completion

const bashSelectCase = `
      if [[ "${cur}" == *","* ]]; then
        local prefix="${cur%%,*},"
        local suffix="${cur##*,}"
        local completions=($(compgen -W "%s" -- "${suffix}"))
        COMPREPLY=()
        for comp in "${completions[@]}"; do
          COMPREPLY+=("${prefix}${comp}")
        done
      else
        COMPREPLY=($(compgen -W "%s" -- "${cur}"))
      fi
      compopt -o nospace 2>/dev/null || true
      return 0
      ;;`

const bashWhereCase = `
      if [[ "${cur}" != *"=" ]] && [[ "${COMPREPLY[0]}" == *"=" ]]; then
        compopt -o nospace 2>/dev/null || true
      fi
      return 0
      ;;`

const bashDefaultCase = `
      return 0
      ;;`

const zshSelectCase = `      %s | %s)
        local cur="${words[CURRENT]}"
        if [[ "$cur" == *","* ]]; then
          local prefix="${cur%%,*},"
          local suffix="${cur##*,}"
          local -a completions
          completions=(%s)
          local -a matching
          for comp in "${completions[@]//\'/}"; do
            if [[ "$comp" == "$suffix"* ]]; then
              matching+=("$comp")
            fi
          done
          compadd -S '' -p "$prefix" -a matching
        else
          local -a opts
          opts=(%s)
          _describe -t %s '%s' opts -S ''
        fi`

const zshWhereCase = ` -S ''`

const zshCaseSuffix = `
            ;;`

const bashScriptTemplate = `
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
`

const zshScriptTemplate = `
  #compdef qp
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

  _qp_completion "$@"
`
