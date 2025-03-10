package cmd

import "github.com/spf13/cobra"

var printCmd = &cobra.Command{
	Use:     "print",
	Short:   "View the blockchain",
	Long:    "🥽 into the blockchain",
	Example: "block print <chain|block>",
}
