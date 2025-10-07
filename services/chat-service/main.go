package main

import (
	"context"
	"log"
	"net"

	"github.com/wutthichod/sa-connext/services/chat-service/database"
	"github.com/wutthichod/sa-connext/shared/messaging"
	pb "github.com/wutthichod/sa-connext/shared/proto/chat"
	"google.golang.org/grpc"
)

var (
	GrpcAddr    = ":9093"
	mongoURI    = ""
	rabbitMqURI = "amqp://guest:guest@localhost:5672/"
)

func main() {
	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	mongoStore := database.NewMongoDB(ctx, mongoURI)
	db := mongoStore.DB()

	// RabbitMQ connection
	rmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rmq.Close()

	// Setup exchange + queue for gateway
	_, err = rmq.SetupQueue("gateway_chat", "chat", "direct", "chat.gateway", true, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Start gRPC server
	chatServer := grpc.NewServer()
	chatService := NewChatService(db, rmq)
	pb.RegisterChatServiceServer(chatServer, chatService)

	log.Println("Server listening on ", GrpcAddr)
	if err := chatServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
