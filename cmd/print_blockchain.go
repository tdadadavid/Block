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
	Args:    cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.ToLower(args[0])
		value := ""
		if len(args) > 1 {
			value = strings.ToLower(args[1])
		}

		if input == "" {
			logger.Error("empty or wrong input passed", slog.String("got", input))
			os.Exit(100)
		}

		command, ok := CommandToHandlers[input]
		if !ok {
			fmt.Println("not implemented")
			return
		}
		command(value)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	CommandToHandlers["b"] = printBlock
	CommandToHandlers["bc"] = printChain
	CommandToHandlers["last"] = printLastBlockOnChain

}
