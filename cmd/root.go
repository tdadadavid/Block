package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "block",
	Short: "block is a tool for interacting with blockchain",
	Long:  `block is a tool for interacting with blockchain.`,
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
	addBlockCmd.PersistentFlags().String("block", "", "Add block to the blockchain")
	printCmd.PersistentFlags().String("chain", "", "Print the chain information")
	printCmd.PersistentFlags().String("block", "", "Print block information")

	// add commands
	rootCmd.AddCommand(addBlockCmd)
	rootCmd.AddCommand(printCmd)
}

//TODO: improve the flag system
