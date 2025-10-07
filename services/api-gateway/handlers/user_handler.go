package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wutthichod/sa-connext/services/api-gateway/grpc_clients/user_client"
	"github.com/wutthichod/sa-connext/services/api-gateway/models"
	pb "github.com/wutthichod/sa-connext/shared/proto/user"
)

type UserHandler struct {
	UserClient *user_client.UserServiceClient
}

func NewUserHandler(uc *user_client.UserServiceClient) *UserHandler {
	return &UserHandler{UserClient: uc}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// Parse incoming JSON
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username and password are required")
	}

	// Call gRPC CreateUser
	resp, err := h.UserClient.CreateUser(c.Context(), &pb.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Contact: &pb.Contact{
			Email: req.Contact.Email,
			Phone: req.Contact.Phone,
		},
		Education: &pb.Education{
			University: req.Education.University,
			Major:      req.Education.Major,
		},
		JobTitle:  req.JobTitle,
		Interests: req.Interests,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return gRPC response to HTTP client
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": resp.GetSuccess(),
	})
}
