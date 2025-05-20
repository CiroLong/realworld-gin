package service

// user_service.go implements UserService
// It handles password hashing, JWT creation, validation logic

import (
	"errors"
	"github/CiroLong/realworld-gin/internal/models"
	"github/CiroLong/realworld-gin/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SaveOneUser(m *models.UserModel) error

	CheckPassword(user *models.UserModel, password string) error
	SetPassword(user *models.UserModel, password string) error

	FindOneUser(userCondition *models.UserModel) (*models.UserModel, error)

	// Register(ctx context.Context, req *models.RegisterRequest) (*model.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepository()}
}

func (us *userService) SaveOneUser(m *models.UserModel) error {
	return us.userRepo.Save(m)
}

func (us *userService) SetPassword(user *models.UserModel, password string) error {
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

func (us *userService) CheckPassword(user *models.UserModel, password string) error {

	if user == nil {
		return errors.New("user not exist")
	}
	bytePassword := []byte(password)
	byteHashedPassword := []byte(user.PasswordHash)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	return err
}

func (us *userService) FindOneUser(userCondition *models.UserModel) (*models.UserModel, error) {
	return us.userRepo.FindOneUser(userCondition)
}
