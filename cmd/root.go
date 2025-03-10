package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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

func init() {
	addBlockCmd.PersistentFlags().String("block", "", "Add block to the blockchain")
	printCmd.PersistentFlags().String("chain", "", "Print the chain information")
	printCmd.PersistentFlags().String("block", "", "Print block information")

	// add commands
	rootCmd.AddCommand(addBlockCmd)
	rootCmd.AddCommand(printCmd)
}

//TODO: Flags are not working and badger db inter-faringing
