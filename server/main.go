package main

import (
	"context"
	"flag"
	"fmt"
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

	topics         map[string][]chan *pb.Message
	messageStorage map[string][]*pb.Message
	mu             sync.Mutex
}

func (b *broker) storeMessage(req *pb.PublishRequest) {
	// 假设我们将消息存储到某个数据库或内存结构中，并在30秒后删除
	b.mu.Lock()
	if _, ok := b.messageStorage[req.Topic]; !ok {
		b.messageStorage[req.Topic] = make([]*pb.Message, 0)
	}
	b.messageStorage[req.Topic] = append(b.messageStorage[req.Topic], req.Message)
	b.mu.Unlock()

	time.AfterFunc(10*time.Second, func() {
		b.mu.Lock()
		b.messageStorage[req.Topic] = b.messageStorage[req.Topic][1:]
		b.mu.Unlock()
	})
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

	go b.storeMessage(req)

	return &pb.PublishResponse{Success: true}, nil
}

func (b *broker) Subscribe(req *pb.SubscribeRequest, stream pb.Broker_SubscribeServer) error {
	if req.Topic == "" {
		return fmt.Errorf("empty topic")
	}

	b.mu.Lock()
	subscriberChan := make(chan *pb.Message, 10)
	if _, ok := b.topics[req.Topic]; !ok {
		b.topics[req.Topic] = make([]chan *pb.Message, 0)
	}
	b.topics[req.Topic] = append(b.topics[req.Topic], subscriberChan)

	if _, ok := b.messageStorage[req.Topic]; ok {
		messages := b.messageStorage[req.Topic]
		for _, msg := range messages {
			select {
			case subscriberChan <- msg:
			default:
			}
		}
	}
	b.mu.Unlock()

	timer := time.After(messageExpiryTime)

	for {
		select {
		case msg := <-subscriberChan:
			if err := stream.Send(msg); err != nil {
				return err
			}
		case <-timer:
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
	topics := make(map[string][]chan *pb.Message)
	messageStorage := make(map[string][]*pb.Message)
	pb.RegisterBrokerServer(s, &broker{topics: topics, messageStorage: messageStorage})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
