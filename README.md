# oops

Undo for your terminal. A shell hook that backs up files before destructive commands run and lets you restore them with one command.

## Install

```
curl -fsSL oops-cli.com/install.sh | bash
```

Then add the hook to your shell:

```sh
# zsh (~/.zshrc)
eval "$(oops init zsh)"

# bash (~/.bashrc)
eval "$(oops init bash)"

# fish (~/.config/fish/config.fish)
oops init fish | source
```

Restart your shell. Run `oops doctor` to verify.

## Usage

```
$ rm -rf src/
▲ rm -r ~/project/src

$ oops
✓ Undid: rm -r ~/project/src
↩ restored ~/project/src
```

`oops 2` undoes the second-to-last action. `oops log` shows history.

## Supported commands

| Command | What oops does | Undo |
|---|---|---|
| `rm` / `rm -rf` | Copies files to trash | restore |
| `mv a b` | Backs up overwrite target | restore b |
| `> file.txt` | Snapshots before redirect | restore |
| `sed -i` | Copies before in-place edit | restore |
| `chmod` / `chown` | Records permissions | restore |
| `git reset --hard` | Creates stash | stash apply |
| `git checkout .` | Creates stash | stash apply |
| `git branch -D` | Logs SHA | recreate |
| `git clean -fd` | Stashes untracked files | stash apply |
| `truncate` | Copies file | restore |

## Commands

| Command | Description |
|---|---|
| `oops` | Undo last action (pass a number to go further back) |
| `oops log` | Show undo history |
| `oops size` | Show backup disk usage |
| `oops clean` | Remove old backups |
| `oops config` | View/set configuration |
| `oops doctor` | Health check |
| `oops init <shell>` | Print shell hook |

## How it works

A `preexec` shell hook pattern-matches each command. Non-destructive commands pass through with zero overhead (no subprocess). Destructive commands trigger `oops protect`, which copies affected files to `~/.oops/trash/` in ~10ms, then lets the original command run.

## License

MIT
