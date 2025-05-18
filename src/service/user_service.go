package service

import (
	"errors"
	"github/CiroLong/realworld-gin/src/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
}

func (us *UserService) SetPassword(userid uint, password string) error {
	u, err := models.GetUser(userid)
	if err != nil {
		return err
	}

	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)

	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (us *UserService) CheckPassword(userid uint, password string) error {
	u, err := models.GetUser(userid)
	if err != nil {
		return err
	}
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
