package repository

import (
	"github/CiroLong/realworld-gin/internal/common"
	"github/CiroLong/realworld-gin/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(data interface{}) error

	FindOneUserByID(userId uint) (*model.UserModel, error)
	FindOneUser(condition interface{}) (*model.UserModel, error)
}

type userRepository struct {
	db *gorm.DB
}

func (ur userRepository) Save(data interface{}) error {
	return ur.db.Save(data).Error
}

func NewUserRepository() UserRepository {
	return &userRepository{db: common.GetDB()}
}

func (ur userRepository) FindOneUserByID(userId uint) (*model.UserModel, error) {
	user := &model.UserModel{}
	err := ur.db.First(user, userId).Error
	return user, err
}

// FindOneUser
// userModel, err := FindOneUser(&UserModel{Username: "username0"})
func (ur userRepository) FindOneUser(condition interface{}) (*model.UserModel, error) {
	model := &model.UserModel{}
	err := ur.db.Where(condition).First(&model).Error
	return model, err
}
