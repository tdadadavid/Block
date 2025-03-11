package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var CommandToHandlers =  map[string]func(){}

var (
	logger *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   "block",
	Short: "🧱🝙",
	Long:  `block is a tool for interacting with chain.`,
	Run: func(cmd *cobra.Command, args []string) {
		data := args[0]
		fmt.Println(data)
	},
	Example: "block add b <DATA>",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	printCmd.PersistentFlags().String("chain", "", "Print the chain information")

	// initialize logger for project
	InitLogger()
}

func InitLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
