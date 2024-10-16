package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	pb "github.com/golrice/pubsub/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

var topic string
var message string

var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"pub", "p"},
	Short:   "Publish a message to a topic",
	Long:    `Publish a message to a topic`,
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if topic != "" && message != "" {
			if err := publish(topic, message); err != nil {
				fmt.Println("publish failed: ", err)
				return
			}
			fmt.Println("you publish", message, "to", topic, "successfully")
		}
	},
}

func init() {
	publishCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic to publish to")
	publishCmd.Flags().StringVarP(&message, "message", "m", "", "Message to publish")
	rootCmd.AddCommand(publishCmd)
}

func publish(topic string, msg string) error {
	flag.Parse()
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := []byte(msg)
	req := &pb.PublishRequest{
		Topic:   topic,
		Message: &pb.Message{Data: message},
	}

	response, err := client.Publish(ctx, req)
	if err != nil {
		return err
	}
	if !response.Success {
		return errors.New("publish failed")
	}

	return nil
}
