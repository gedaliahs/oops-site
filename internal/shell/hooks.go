package shell

import "fmt"

// ZshHook returns the zsh preexec hook script.
func ZshHook(oopsBin string) string {
	return fmt.Sprintf(`# oops — terminal undo (zsh hook)
_oops_preexec() {
  local cmd="$1"
  # Fast path: only invoke oops for potentially destructive commands
  case "$cmd" in
    rm\ *|rm|mv\ *|sed\ *|gsed\ *|chmod\ *|chown\ *|truncate\ *|gtruncate\ *|git\ reset*|git\ checkout*|git\ clean*|git\ branch\ *-[dD]*)
      %s protect -- "$cmd"
      ;;
    *\>*)
      # Redirect — check if overwriting an existing file
      %s protect-redirect -- "$cmd"
      ;;
  esac
}
autoload -Uz add-zsh-hook
add-zsh-hook preexec _oops_preexec
`, oopsBin, oopsBin)
}

// BashHook returns the bash DEBUG trap hook script.
func BashHook(oopsBin string) string {
	return fmt.Sprintf(`# oops — terminal undo (bash hook)
_oops_preexec() {
  # Only run on real commands, not prompts
  [ -n "$COMP_LINE" ] && return
  [ "$BASH_COMMAND" = "$PROMPT_COMMAND" ] && return

  local cmd="$BASH_COMMAND"
  case "$cmd" in
    rm\ *|rm|mv\ *|sed\ *|gsed\ *|chmod\ *|chown\ *|truncate\ *|gtruncate\ *|git\ reset*|git\ checkout*|git\ clean*|git\ branch\ *-[dD]*)
      %s protect -- "$cmd"
      ;;
    *\>*)
      %s protect-redirect -- "$cmd"
      ;;
  esac
}
trap '_oops_preexec' DEBUG
`, oopsBin, oopsBin)
}

// FishHook returns the fish preexec hook script.
func FishHook(oopsBin string) string {
	return fmt.Sprintf(`# oops — terminal undo (fish hook)
function _oops_preexec --on-event fish_preexec
  set -l cmd $argv[1]
  switch $cmd
    case 'rm *' 'mv *' 'sed *' 'gsed *' 'chmod *' 'chown *' 'truncate *' 'gtruncate *' 'git reset*' 'git checkout*' 'git clean*' 'git branch *-D*' 'git branch *-d*'
      %s protect -- "$cmd"
    case '*>*'
      %s protect-redirect -- "$cmd"
  end
end
`, oopsBin, oopsBin)
}
