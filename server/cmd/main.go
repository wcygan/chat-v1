package main

import (
	"context"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	pb "github.com/wcygan/chat-v1/generated/go/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

func main() {
	// Start an embedded NATS server
	opts := &server.Options{
		Host: "localhost",
		Port: 4222,
	}
	natsServer, err := server.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	go natsServer.Start()
	if !natsServer.ReadyForConnections(10 * time.Second) {
		log.Fatal("NATS server did not start in time")
	} else {
		log.Println("NATS server is ready")

	}

	// Connect to the NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("NATS client is connected")
	}
	defer nc.Close()

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Println("gRPC server is ready and serving at :50051")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &chatServer{nc: nc})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type chatServer struct {
	pb.UnimplementedChatServiceServer
	nc *nats.Conn
}

func (s *chatServer) JoinChat(req *pb.JoinChatRequest, stream pb.ChatService_JoinChatServer) error {
	// Create a channel to signal when the stream is closed
	done := make(chan struct{})

	// Subscribe to the chat room topic
	sub, err := s.nc.Subscribe(req.ChatRoom, func(m *nats.Msg) {
		var msg pb.ChatMessage
		// Deserialize the protobuf message
		if err := proto.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("Error deserializing message: %v", err)
			return
		}
		log.Printf("Sending message from %s to %s in chat room %s: %s (UUID: %s)", msg.User, req.User, msg.ChatRoom, msg.Message, msg.Uuid) // Update this line
		if err := stream.Send(&msg); err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
	})
	if err != nil {
		return err
	}

	// Wait for the client to close the stream
	<-stream.Context().Done()
	log.Printf("%s disconnected from chat room %s", req.User, req.ChatRoom)

	// Unsubscribe from the NATS topic
	if err := sub.Unsubscribe(); err != nil {
		log.Printf("Error unsubscribing: %v", err)
	}

	// Signal that the stream is closed
	close(done)

	return nil
}

func (s *chatServer) SendChatMessage(ctx context.Context, msg *pb.ChatMessage) (*emptypb.Empty, error) {
	log.Printf("Received message from %s in chat room %s: %s (UUID: %s)", msg.User, msg.ChatRoom, msg.Message, msg.Uuid) // Update this line
	// Serialize the protobuf message
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// Publish the serialized message to the chat room topic
	err = s.nc.Publish(msg.ChatRoom, data)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
