package repository

import (
	"context"

	"github.com/wutthichod/sa-connext/services/user-service/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
