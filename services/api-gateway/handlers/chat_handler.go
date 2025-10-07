package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wutthichod/sa-connext/services/api-gateway/grpc_clients/chat_client"
	"github.com/wutthichod/sa-connext/services/api-gateway/models"
	"github.com/wutthichod/sa-connext/shared/messaging"
	pb "github.com/wutthichod/sa-connext/shared/proto/chat"
)

type ChatHandler struct {
	ChatClient  *chat_client.ChatServiceClient
	ConnManager *messaging.ConnectionManager
	Queue       *messaging.QueueConsumer // listens to messages from chat service
}

func NewChatHandler(client *chat_client.ChatServiceClient, connManager *messaging.ConnectionManager, queue *messaging.QueueConsumer) *ChatHandler {
	return &ChatHandler{ChatClient: client, ConnManager: connManager, Queue: queue}
}

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	var req models.CreateChatRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	if req.UserID == "" || req.RecipientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "UserID and RecipientID are required")
	}

	res, err := h.ChatClient.CreateChat(c.Context(), &pb.CreateChatRequest{
		UserId:      req.UserID,
		RecipientId: req.RecipientID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	var req models.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	if req.UserID == "" || req.RecipientID == "" || req.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "UserID, RecipientID, and Message are required")
	}

	// Call ChatService via gRPC
	_, err := h.ChatClient.SendMessage(c.Context(), &pb.SendMessageRequest{
		FromUserId: req.UserID,
		ToUserId:   req.RecipientID,
		Message:    req.Message,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// ListenRabbit runs in a background goroutine to relay RabbitMQ messages to WebSocket clients
func (h *ChatHandler) ListenRabbit() {
	err := h.Queue.Start()
	if err != nil {
		log.Fatal("Failed to start RabbitMQ consumer:", err)
	}
}
