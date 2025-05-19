package repository

import (
	"context"
	"github/CiroLong/realworld-gin/internal/common"
	"github/CiroLong/realworld-gin/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.UserModel) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: common.GetDB()}
}

func (u userRepository) CreateUser(ctx context.Context, user *models.UserModel) error {
	//TODO implement me
	panic("implement me")
}
