# oops

Undo for your terminal. A shell hook that backs up files before destructive commands run and lets you restore them with one command.

## Install

```
curl -fsSL oops-cli.com/install.sh | bash
```

The installer handles everything — downloads the binary, adds the shell hook to your shell config, and creates the backup directory.

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

## Commands

| Command | Description |
|---|---|
| `oops` | Undo last action (pass a number to go further back) |
| `oops log` | Show undo history |
| `oops size` | Show backup disk usage |
| `oops clean` | Remove old backups (`--all` for everything) |
| `oops config` | View or change settings |
| `oops doctor` | Check installation health |
| `oops tutorial` | Interactive walkthrough |
| `oops uninstall` | Remove oops from your system |
| `oops --version` | Print version |
| `oops --upgrade` | Upgrade to the latest version |

## How it works

A `preexec` shell hook pattern-matches each command. Non-destructive commands pass through with zero overhead (no subprocess). Destructive commands trigger `oops protect`, which copies affected files to `~/.oops/trash/` in ~10ms, then lets the original command run.

## Uninstall

```
oops uninstall
```

Removes the shell hook and backup directory. Then run `sudo rm /usr/local/bin/oops` to remove the binary.

## License

MIT
