package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addBlockCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"addition"},
	Short:   "Use to add a new block",
	Long:    "What is a chain without a block 🧱",
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.ToLower(args[0]) // everything on the chain is converted to small letters
		if input == "" {
			logger.Error("empty or wrong input passed", slog.String("expected", "<anything>"), slog.String("got", ""))
			os.Exit(100)
		}
		addBlock(input)
	},
	Args:    cobra.ExactArgs(1),
	Example: "block add <DATA>",
}

func init() {
	rootCmd.AddCommand(addBlockCmd)
}
