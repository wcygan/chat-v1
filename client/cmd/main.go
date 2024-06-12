package main

import (
	"context"
	"fmt"
	pb "github.com/wcygan/chat-v1/generated/go/chat/v1"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Println("gRPC connection established")
	defer conn.Close()
	c := pb.NewChatServiceClient(conn)

	// Join chat
	joinCtx, joinCancel := context.WithCancel(context.Background())
	defer joinCancel()
	stream, err := c.JoinChat(joinCtx, &pb.JoinChatRequest{
		User:     "user1",
		ChatRoom: "room1",
	})
	if err != nil {
		log.Fatalf("could not join chat: %v", err)
	}
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Println("Stream closed by server")
				break // Exit the loop if stream is closed
			} else if err != nil {
				log.Fatalf("Failed to receive a message: %v", err)
			}
			fmt.Printf("Received message %s from %s\n", in.Message, in.User)
		}
	}()

	// Send a message
	sendCtx, sendCancel := context.WithCancel(context.Background())
	defer sendCancel()
	_, err = c.SendChatMessage(sendCtx, &pb.ChatMessage{
		User:     "user1",
		ChatRoom: "room1",
		Message:  "Hello, World!",
	})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	// Keep the client running to listen for incoming messages
	select {}
}
