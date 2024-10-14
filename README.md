# Pub-Sub System using Protobuf and gRPC

## Overview

This project implements a **Publish-Subscribe (Pub-Sub) system** using **Protobuf** and **gRPC**. It allows publishers to send messages to a specific topic, and subscribers to receive messages from the topics they are subscribed to. The system supports message persistence with an expiration time on the server side, after which the messages are deleted.

## Features

- **Publisher**: Send messages to specific topics.
- **Subscriber**: Subscribe to topics and receive real-time messages.
- **Broker**: Manage message distribution and handle topic subscriptions.
- **Message Expiration**: Messages on the server are stored temporarily and are removed after a configurable time.
- **gRPC Communication**: Provides efficient and scalable communication between services using gRPC.
  
## Architecture

The project follows a microservices-based architecture, consisting of three main components:

1. **Publisher**: Clients that publish messages to topics.
2. **Subscriber**: Clients that subscribe to topics and receive messages in real time.
3. **Broker**: The central server that handles message publishing, topic subscription, and message expiration.

### System Architecture Diagram

```
             +------------------+  
             |   Publisher       |  
             +--------+----------+  
                      |  gRPC  
                      v  
             +------------------+  
             |     Broker        |  <------- (Message Storage with Expiry)
             +--------+----------+  
                      |  gRPC  
                      v  
             +------------------+  
             |   Subscriber      |  
             +------------------+  
```

## Getting Started

### Prerequisites

Make sure you have the following installed on your machine:

- [Go](https://golang.org/dl/) (>=1.18)
- [Protobuf](https://developers.google.com/protocol-buffers) Compiler (protoc)
- [gRPC-Go](https://grpc.io/docs/languages/go/quickstart/) library
- [Docker](https://www.docker.com/) (optional for deployment)
  
### Installation

1. Clone the repository:

   ```bash
   git clone git@github.com:golrice/pubsub.git
   cd pubsub
   ```

2. Install the required Go packages:

   ```bash
   go mod tidy
   ```

3. Compile Protobuf definitions:

   ```bash
   protoc --go_out=. --go-grpc_out=. proto/*.proto
   ```

### Running the System

1. **Start the Broker (Server)**:

   ```bash
   go run server/main.go
   ```

2. **Run a Publisher**:

   ```bash
   go run client/publisher/main.go --topic=example --message="Hello World!"
   ```

3. **Run a Subscriber**:

   ```bash
   go run client/subscriber/main.go --topic=example
   ```

### Configuration

You can configure various aspects of the system, such as:

- **Message expiration time**: Modify the message expiration time in the `server/config.go` file.
  
## Project Structure

```
pubsub/
├── client/
│   ├── publisher/            # Publisher client implementation
│   ├── subscriber/           # Subscriber client implementation
│   └── utils/                # Utility functions for client operations
├── proto/
│   ├── pubsub.proto          # Protobuf definitions
├── server/
│   ├── main.go               # Broker (server) implementation
│   ├── config.go             # Configuration settings
│   └── storage.go            # Message storage and expiration logic
├── go.mod                    # Go module file
├── go.sum                    # Dependency file
└── README.md                 # Project README
```

## How It Works

1. **Publisher** sends a message to a specified topic via the `Publish` gRPC method.
2. **Broker** processes the message and stores it temporarily. It also forwards the message to all subscribers of the topic.
3. **Subscriber** listens for messages in real-time via the `Subscribe` gRPC method and receives messages as long as the subscription is active.
4. **Message Expiration**: After the configured expiration time, the broker deletes the message from its storage.

## Example Usage

### Publishing a Message

To publish a message to a topic:

```bash
go run client/publisher/main.go --topic=example --message="This is a test message."
```

### Subscribing to a Topic

To subscribe to a topic:

```bash
go run client/subscriber/main.go --topic=example
```

### Configuring Message Expiration

In `server/config.go`, you can set the expiration time for messages:

```go
const messageExpiryTime = 30 * time.Second  // Set message expiration to 30 seconds
```

## Future Enhancements

- **Message Persistence**: Store messages in a database for durability and historical lookup.
- **Load Balancing**: Implement a distributed broker with load balancing for better scalability.
- **Authorization**: Add authentication and authorization mechanisms for secure communication.

## License

This project is licensed under the MIT License.

## Contributing

Feel free to submit issues and pull requests for new features, bug fixes, and improvements.
