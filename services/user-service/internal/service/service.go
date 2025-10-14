package service

import (
	"context"
	"log"

	"github.com/wutthichod/sa-connext/services/user-service/internal/mapper"
	"github.com/wutthichod/sa-connext/services/user-service/internal/repository"
	"github.com/wutthichod/sa-connext/shared/contracts"
	"github.com/wutthichod/sa-connext/shared/messaging"
	pb "github.com/wutthichod/sa-connext/shared/proto/user"
)

type Service interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) error
}

type service struct {
	repo repository.Repository
	rb   *messaging.RabbitMQ
}

func NewService(repo repository.Repository, rb *messaging.RabbitMQ) Service {
	return &service{repo, rb}
}

func (s *service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) error {
	// PB → DTO
	dtoUser := mapper.FromPbRequest(req)
	log.Printf("Mapped DTO: %+v\n", dtoUser)
	// DTO → Model
	userModel := mapper.ToUserModel(dtoUser)

	log.Printf("Mapped Model: %+v\n", userModel)

	// Publish to RabbitMQ
	event := contracts.EmailEvent{
	To:      "brightka.ceo@gmail.com",
	Subject: "Welcome!",
	Body:    "Hi there, thanks for signing up!",
	}
	
	if err := s.rb.PublishMessage(context.Background(),"notification.exchange","notification.email",event); err != nil {
			log.Printf("Failed to publish email event: %v", err)
	}

	// Save to DB
	return s.repo.CreateUser(ctx, userModel)
}
