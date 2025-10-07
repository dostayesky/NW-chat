package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wutthichod/sa-connext/shared/messaging"
	pb "github.com/wutthichod/sa-connext/shared/proto/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	db  *mongo.Database
	rmq *messaging.RabbitMQ
}

func NewChatService(db *mongo.Database, rmq *messaging.RabbitMQ) *ChatService {
	return &ChatService{db: db, rmq: rmq}
}

// CreateChat creates a new chat between two users
func (s *ChatService) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	chat := bson.M{
		"user_ids":   []string{req.UserId, req.RecipientId},
		"created_at": time.Now(),
	}

	res, err := s.db.Collection("chat").InsertOne(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	chatID := res.InsertedID.(primitive.ObjectID).Hex()

	return &pb.CreateChatResponse{
		UserId:      req.UserId,
		RecipientId: req.RecipientId,
		ChatId:      chatID,
	}, nil
}

// SendMessage saves a message and publishes to RabbitMQ
func (s *ChatService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Find the chat document for the two users
	var chat bson.M
	err := s.db.Collection("chat").FindOne(ctx, bson.M{
		"user_ids": bson.M{"$all": []string{req.FromUserId, req.ToUserId}},
	}).Decode(&chat)
	if err != nil {
		return nil, fmt.Errorf("chat not found: %v", err)
	}

	chatID := chat["_id"].(primitive.ObjectID)

	// Insert the message
	messageDoc := bson.M{
		"chat_id":      chatID,
		"from_user_id": req.FromUserId,
		"to_user_id":   req.ToUserId,
		"content":      req.Message,
		"created_at":   time.Now(),
	}

	msgRes, err := s.db.Collection("messages").InsertOne(ctx, messageDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to save message: %v", err)
	}

	messageID := msgRes.InsertedID.(primitive.ObjectID).Hex()

	// Publish to RabbitMQ
	msg := map[string]interface{}{
		"chat_id":      chatID.Hex(),
		"message_id":   messageID,
		"from_user_id": req.FromUserId,
		"to_user_id":   req.ToUserId,
		"message":      req.Message,
	}

	if err := s.rmq.PublishMessage(ctx, "chat", "chat.gateway", msg); err != nil {
		log.Printf("failed to publish message to RabbitMQ: %v", err)
	}

	return &pb.SendMessageResponse{
		MessageId: messageID,
		Status:    "sent",
	}, nil
}
