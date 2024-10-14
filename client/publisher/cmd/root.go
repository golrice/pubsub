package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "publisher",
	Short: "A command line tool for publishing messages to a pubsub topic",
	Long:  "A command line tool for publishing messages to a pubsub topic",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
