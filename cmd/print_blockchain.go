package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:     "print",
	Short:   "View the chain",
	Long:    "🥽 into the chain",
	Example: "block print <chain|block>",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.ToLower(args[0])

		if input == "" {
			logger.Error("empty or wrong input passed", slog.String("got", input))
			os.Exit(100)
		}

		if CommandToHandlers[input] == nil {
			fmt.Println("Not implemented")
		}

		command := CommandToHandlers[input]
		command()
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	CommandToHandlers["bc"] = printChain
	CommandToHandlers["last"] = printLastBlockOnChain
}
