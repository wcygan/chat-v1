package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	pb "github.com/wcygan/chat-v1/generated/go/chat/v1"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

var (
	username   string
	chatroom   string
	clientUUID string
)

func main() {
	clientUUID = uuid.New().String()

	var rootCmd = &cobra.Command{
		Use:   "chat-client",
		Short: "Chat client to join and send messages to a chat room",
		Run: func(cmd *cobra.Command, args []string) {
			runClient()
		},
	}

	rootCmd.Flags().StringVarP(&username, "username", "u", "", "Username for the chat")
	rootCmd.Flags().StringVarP(&chatroom, "chatroom", "c", "", "Chatroom to join")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("chatroom")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatServiceClient(conn)

	// Join chat
	joinCtx, joinCancel := context.WithCancel(context.Background())
	defer joinCancel()
	stream, err := c.JoinChat(joinCtx, &pb.JoinChatRequest{
		User:     username,
		ChatRoom: chatroom,
	})
	if err != nil {
		log.Fatalf("could not join chat: %v", err)
	} else {
		log.Printf("Joined chatroom %s as %s", chatroom, username)
	}

	// Goroutine to receive messages
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Println("Stream closed by server")
				break
			} else if err != nil {
				log.Fatalf("Failed to receive a message: %v", err)
			} else if in.Uuid == clientUUID {
				// Skip messages sent by this client
			} else {
				fmt.Printf("[%s]  %s\n", in.User, in.Message)
			}
		}
	}()

	// Allow user to send messages interactively
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		sendCtx, sendCancel := context.WithCancel(context.Background())
		defer sendCancel()
		_, err = c.SendChatMessage(sendCtx, &pb.ChatMessage{
			User:     username,
			ChatRoom: chatroom,
			Message:  text,
			Uuid:     clientUUID,
		})
		if err != nil {
			log.Fatalf("could not send message: %v", err)
		}
	}
}
