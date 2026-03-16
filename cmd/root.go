package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/gedaliah/oops/internal/cleanup"
	"github.com/gedaliah/oops/internal/journal"
	"github.com/gedaliah/oops/internal/style"
	"github.com/gedaliah/oops/internal/trash"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "oops [N]",
	Short: "Terminal undo — restore your last destructive command",
	Long:  style.Banner() + "\n\nUndo destructive terminal commands. Run `oops` to undo the last action, or `oops N` to undo the Nth most recent.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUndo,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runUndo(cmd *cobra.Command, args []string) error {
	// Lazy cleanup
	cleanup.RunIfNeeded()

	n := 1
	if len(args) == 1 {
		var err error
		n, err = strconv.Atoi(args[0])
		if err != nil || n < 1 {
			return fmt.Errorf("invalid argument: %s (expected a positive number)", args[0])
		}
	}

	entries, err := journal.Last(n)
	if err != nil {
		return fmt.Errorf("reading journal: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println(style.Warning("Nothing to undo"))
		return nil
	}

	if n > len(entries) {
		return fmt.Errorf("only %d undoable actions in history", len(entries))
	}

	entry := entries[n-1]

	// Handle git-specific undos
	if entry.GitAction != "" {
		return undoGit(entry)
	}

	// Standard file restore
	if entry.TrashDir == "" {
		return fmt.Errorf("no backup found for this action")
	}

	restored, err := trash.Restore(entry.TrashDir)
	if err != nil {
		return fmt.Errorf("restoring files: %w", err)
	}

	// Mark as undone
	if err := journal.MarkUndone(entry.ID); err != nil {
		fmt.Fprintln(os.Stderr, style.Warning("Could not mark entry as undone: "+err.Error()))
	}

	fmt.Println(style.Success(fmt.Sprintf("Undid: %s", entry.Desc)))
	for _, f := range restored {
		fmt.Println(style.Restored(f))
	}

	return nil
}

func undoGit(entry journal.Entry) error {
	switch entry.GitAction {
	case "stash":
		stashRef := entry.GitStash
		if stashRef == "" {
			stashRef = "stash@{0}"
		}
		out, err := exec.Command("git", "stash", "apply", stashRef).CombinedOutput()
		if err != nil {
			return fmt.Errorf("git stash apply failed: %s", string(out))
		}
		if err := journal.MarkUndone(entry.ID); err != nil {
			fmt.Fprintln(os.Stderr, style.Warning("Could not mark entry as undone: "+err.Error()))
		}
		fmt.Println(style.Success(fmt.Sprintf("Undid: %s", entry.Desc)))
		fmt.Println(style.Green.Render("  Applied stash: ") + stashRef)

	case "log-branch":
		if entry.GitSHA == "" {
			return fmt.Errorf("no SHA recorded for deleted branch")
		}
		branchName := entry.GitRef
		out, err := exec.Command("git", "branch", branchName, entry.GitSHA).CombinedOutput()
		if err != nil {
			return fmt.Errorf("git branch restore failed: %s", string(out))
		}
		if err := journal.MarkUndone(entry.ID); err != nil {
			fmt.Fprintln(os.Stderr, style.Warning("Could not mark entry as undone: "+err.Error()))
		}
		fmt.Println(style.Success(fmt.Sprintf("Undid: %s", entry.Desc)))
		fmt.Println(style.Green.Render("  Restored branch: ") + branchName + " at " + entry.GitSHA[:8])

	default:
		return fmt.Errorf("unknown git action: %s", entry.GitAction)
	}

	return nil
}
