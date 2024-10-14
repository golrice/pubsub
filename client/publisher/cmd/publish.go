package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var topic string

var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"pub", "p"},
	Short:   "Publish a message to a topic",
	Long:    `Publish a message to a topic`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("publish msg: ", args[0])
		if topic != "" {
			fmt.Println("publish to topic: ", topic)
		}
	},
}

func init() {
	publishCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic to publish to")
	rootCmd.AddCommand(publishCmd)
}
