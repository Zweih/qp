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
