package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/golrice/pubsub/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type broker struct {
	pb.UnimplementedBrokerServer

	topics map[string][]chan *pb.Message
	mu     sync.Mutex
}

func (b *broker) Publish(cxt context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	if req.Topic == "" {
		return nil, fmt.Errorf("empty topic")
	}
	if req.Message == nil {
		return nil, fmt.Errorf("empty message")
	}

	b.mu.Lock()
	for _, subscribers := range b.topics[req.Topic] {
		select {
		case subscribers <- req.Message:
		default:
		}
	}
	b.mu.Unlock()

	return &pb.PublishResponse{Success: true}, nil
}

func (b *broker) Subscribe(req *pb.SubscribeRequest, stream pb.Broker_SubscribeServer) error {
	if req.Topic == "" {
		return fmt.Errorf("empty topic")
	}

	b.mu.Lock()
	subscriberChan := make(chan *pb.Message)
	b.topics[req.Topic] = append(b.topics[req.Topic], subscriberChan)
	b.mu.Unlock()

	for {
		select {
		case msg := <-subscriberChan:
			if err := stream.Send(msg); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	log.Println("Starting server on port", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &broker{topics: make(map[string][]chan *pb.Message)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
