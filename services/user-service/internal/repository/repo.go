package repository

import (
	"context"

	"github.com/wutthichod/sa-connext/services/user-service/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).
		Joins("JOIN contacts ON contacts.id = users.contact_id").
		Where("contacts.email = ?", email).
		Preload("Contact").
		Preload("Education").
		Preload("Interests").
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // not found, return nil user
		}
		return nil, err // DB or query error
	}
	return &user, nil
}
