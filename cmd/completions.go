package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionsCmd = &cobra.Command{
	Use:       "completions <shell>",
	Short:     "Generate shell completions",
	Long:      "Generate shell completion scripts for zsh, bash, or fish.",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"zsh", "bash", "fish"},
	RunE:      runCompletions,
}

func init() {
	rootCmd.AddCommand(completionsCmd)
}

func runCompletions(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	default:
		return fmt.Errorf("unsupported shell: %s", args[0])
	}
}
