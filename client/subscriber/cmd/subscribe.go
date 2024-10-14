package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:     "subscribe",
	Aliases: []string{"sub", "s"},
	Short:   "Subscribe to a topic",
	Long:    "Subscribe to a topic",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe topic: ", args[0])
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
