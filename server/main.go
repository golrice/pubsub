package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/golrice/pubsub/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type broker struct {
	pb.UnimplementedBrokerServer

	topics map[string]chan *pb.Message
	mu     sync.Mutex
}

func Publish(client pb.BrokerClient, topic string, message *pb.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.Publish(ctx, &pb.PublishRequest{Topic: topic, Message: message})
	if err != nil {
		return err
	}

	if response.Success != true {
		return errors.New("publish failed")
	}

	return nil
}

func Subscribe(client pb.BrokerClient, topic string) (chan *pb.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{Topic: topic})
	if err != nil {
		return nil, err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		fmt.Println(string(response.Data))
	}

	return nil, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	log.Println("Starting server on port", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &broker{topics: make(map[string]chan *pb.Message)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
