package service

// user_service.go implements UserService
// It handles password hashing, JWT creation, validation logic

import (
	"errors"
	"github/CiroLong/realworld-gin/internal/model"
	"github/CiroLong/realworld-gin/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SaveOneUser(m *model.UserModel) error

	CheckPassword(user *model.UserModel, password string) error
	SetPassword(user *model.UserModel, password string) error

	FindOneUser(userCondition *model.UserModel) (*model.UserModel, error)

	// Register(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepository()}
}

func (us *userService) SaveOneUser(m *model.UserModel) error {
	return us.userRepo.Save(m)
}

func (us *userService) SetPassword(user *model.UserModel, password string) error {
	if user == nil {
		return errors.New("user not exist")
	}
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)

	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	user.PasswordHash = string(passwordHash)
	return nil
}

func (us *userService) CheckPassword(user *model.UserModel, password string) error {

	if user == nil {
		return errors.New("user not exist")
	}
	bytePassword := []byte(password)
	byteHashedPassword := []byte(user.PasswordHash)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	return err
}

func (us *userService) FindOneUser(userCondition *model.UserModel) (*model.UserModel, error) {
	return us.userRepo.FindOneUser(userCondition)
}
