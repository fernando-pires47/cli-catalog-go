package completion

import "fmt"

func BashScript() string {
	return fmt.Sprintf(`# bash completion for cs
_cs_complete() {
  local cur
  cur="${COMP_WORDS[COMP_CWORD]}"
  local line=("${COMP_WORDS[@]:1}")
  COMPREPLY=( $(compgen -W "$(cs __complete "${line[@]}")" -- "$cur") )
}
complete -F _cs_complete cs
`)
}
