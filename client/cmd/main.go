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
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatServiceClient(conn)

	// Join chat
	stream, err := c.JoinChat(context.Background(), &pb.JoinChatRequest{
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
	_, err = c.SendChatMessage(context.Background(), &pb.ChatMessage{
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
