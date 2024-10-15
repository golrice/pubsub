package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"

	pb "github.com/golrice/pubsub/proto"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

var subscribeCmd = &cobra.Command{
	Use:     "subscribe",
	Aliases: []string{"sub", "s"},
	Short:   "Subscribe to a topic",
	Long:    "Subscribe to a topic",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe topic: ", args[0])
		if args[0] != "" {
			if err := Subscribe(args[0]); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}

func Subscribe(topic string) error {
	flag.Parse()
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	ctx := context.Background()

	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{Topic: topic})
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", topic, string(msg.Data))
	}

	return nil
}
