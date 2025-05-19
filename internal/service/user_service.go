package service

// user_service.go implements UserService
// It handles password hashing, JWT creation, validation logic

import (
	"github/CiroLong/realworld-gin/internal/repository"
)

type UserService interface {
	// Register(ctx context.Context, req *models.RegisterRequest) (*model.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepository()}
}

func (us *userService) SetPassword(userid uint, password string) error {
	// TODO:

	//u, err := models.GetUser(userid)
	//if err != nil {
	//	return err
	//}
	//
	//if len(password) == 0 {
	//	return errors.New("password should not be empty!")
	//}
	//bytePassword := []byte(password)
	//
	//passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	//u.PasswordHash = string(passwordHash)
	return nil
}

func (us *userService) CheckPassword(userid uint, password string) error {
	// TODO:

	//u, err := models.GetUser(userid)
	//if err != nil {
	//	return err
	//}
	//bytePassword := []byte(password)
	//byteHashedPassword := []byte(u.PasswordHash)
	//return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	return nil
}
