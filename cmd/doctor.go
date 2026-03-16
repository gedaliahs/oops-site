package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gedaliah/oops/internal/config"
	"github.com/gedaliah/oops/internal/style"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check oops installation health",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println(style.Banner())
	fmt.Println()

	ok := true

	// Check oops directory
	if _, err := os.Stat(config.OopsDir()); err != nil {
		fmt.Println(style.Error("~/.oops directory missing"))
		ok = false
	} else {
		fmt.Println(style.Success("~/.oops directory exists"))
	}

	// Check trash directory
	if _, err := os.Stat(config.TrashDir()); err != nil {
		fmt.Println(style.Error("~/.oops/trash directory missing"))
		ok = false
	} else {
		fmt.Println(style.Success("~/.oops/trash directory exists"))
	}

	// Check config
	cfg := config.Load()
	fmt.Println(style.Success(fmt.Sprintf("Config: retention=%dd, max_trash=%s", cfg.RetentionDays, style.FormatSize(cfg.MaxTrashBytes))))

	// Check journal
	if _, err := os.Stat(config.JournalPath()); err != nil {
		fmt.Println(style.Dim.Render("  " + style.SymBackup + " No journal yet (normal for first run)"))
	} else {
		fmt.Println(style.Success("Journal file exists"))
	}

	// Check git availability
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println(style.Warning("git not found (git-related undo will not work)"))
	} else {
		fmt.Println(style.Success("git available"))
	}

	// Check shell hook
	shellVar := os.Getenv("SHELL")
	if shellVar != "" {
		fmt.Println(style.Success("Shell: " + shellVar))
	}

	// Check if hook is active
	oopsHook := os.Getenv("_OOPS_HOOK")
	if oopsHook == "1" {
		fmt.Println(style.Success("Shell hook is active"))
	} else {
		fmt.Println(style.Warning("Shell hook not detected (add `eval \"$(oops init zsh)\"` to your shell rc)"))
	}

	fmt.Println()
	if ok {
		fmt.Println(style.Success("All checks passed"))
	} else {
		fmt.Println(style.Error("Some checks failed — run `oops init <shell>` for setup instructions"))
	}

	return nil
}
