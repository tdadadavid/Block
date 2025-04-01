package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dropCmd = &cobra.Command{
	Use:     "drop",
	Short:   "clear everything",
	Long:    "This will destroy the blockchain[❌]",
	Example: "block drop",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		err := os.Remove(chainStorePath)
		if err != nil {
			fmt.Printf("failed to destroy chain: %s", err)
			return
		}
	},
}
