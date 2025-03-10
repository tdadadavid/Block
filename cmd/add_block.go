package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var addBlockCmd = &cobra.Command{
	Use:     "AddBlock",
	Short:   "Use to add a new block",
	Long:    "What is a blockchain without a block 🧱",
	Run:     Handle,
	Args:    cobra.ExactArgs(1),
	Example: "block add b <DATA>",
}

func Handle(cmd *cobra.Command, args []string) {
	data := args[0]
	fmt.Println("data:", data)
}
