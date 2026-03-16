package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gedaliah/oops/internal/config"
	"github.com/gedaliah/oops/internal/style"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove oops from your system",
	RunE:  runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	fmt.Println(style.Bold.Render("oops uninstaller"))
	fmt.Println()

	// 1. Remove shell hook from rc file
	home, _ := os.UserHomeDir()
	shell := filepath.Base(os.Getenv("SHELL"))

	rcFiles := map[string]string{
		"zsh":  filepath.Join(home, ".zshrc"),
		"bash": filepath.Join(home, ".bashrc"),
		"fish": filepath.Join(home, ".config", "fish", "config.fish"),
	}

	rcFile := rcFiles[shell]
	if rcFile == "" {
		// Try all of them
		for _, f := range rcFiles {
			removeHookFromFile(f)
		}
	} else {
		removeHookFromFile(rcFile)
	}

	// 2. Remove ~/.oops directory
	oopsDir := config.OopsDir()
	if _, err := os.Stat(oopsDir); err == nil {
		if confirm(fmt.Sprintf("Remove %s and all backups?", oopsDir)) {
			os.RemoveAll(oopsDir)
			fmt.Println(style.Success("Removed " + oopsDir))
		} else {
			fmt.Println(style.Dim.Render("  Skipped"))
		}
	}

	// 3. Remove the binary
	bin, err := os.Executable()
	if err == nil {
		if confirm(fmt.Sprintf("Remove binary %s?", bin)) {
			if err := os.Remove(bin); err != nil {
				fmt.Println(style.Warning("Could not remove binary (try: sudo rm " + bin + ")"))
			} else {
				fmt.Println(style.Success("Removed " + bin))
			}
		} else {
			fmt.Println(style.Dim.Render("  Skipped"))
		}
	}

	fmt.Println()
	fmt.Println(style.Success("oops has been uninstalled. Open a new terminal tab to finish."))
	return nil
}

func removeHookFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	var filtered []string
	found := false
	for _, line := range lines {
		if strings.Contains(line, "oops init") {
			found = true
			continue
		}
		filtered = append(filtered, line)
	}

	if !found {
		return
	}

	if confirm(fmt.Sprintf("Remove shell hook from %s?", path)) {
		// Remove trailing empty lines left behind
		for len(filtered) > 0 && filtered[len(filtered)-1] == "" {
			filtered = filtered[:len(filtered)-1]
		}
		filtered = append(filtered, "") // single trailing newline

		os.WriteFile(path, []byte(strings.Join(filtered, "\n")), 0o644)
		fmt.Println(style.Success("Removed hook from " + path))
	} else {
		fmt.Println(style.Dim.Render("  Skipped"))
	}
}

func confirm(msg string) bool {
	fmt.Printf("  %s [Y/n] ", msg)
	reader := bufio.NewReader(os.Stdin)
	reply, _ := reader.ReadString('\n')
	reply = strings.TrimSpace(reply)
	return reply == "" || strings.HasPrefix(strings.ToLower(reply), "y")
}
