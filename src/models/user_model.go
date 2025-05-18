package models

import (
	"github/CiroLong/realworld-gin/src/common"
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique_index"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
}

type FollowModel struct {
	gorm.Model
	Following    UserModel
	FollowingID  uint
	FollowedBy   UserModel
	FollowedByID uint
}

func GetUser(userId uint) (*UserModel, error) {
	u, err := FindOneUser(&UserModel{ID: userId})
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindOneUser
//
//	You could input the conditions and it will return an UserModel in database with error info.
//
//	userModel, err := FindOneUser(&UserModel{Username: "username0"})
func FindOneUser(condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

// SaveOne
//
//	You could input an UserModel which will be saved in database returning with error info
//
//	if err := SaveOne(&userModel); err != nil { ... }
func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func (u *UserModel) Create() (userId uint, err error) {
	db := common.GetDB()
	result := db.Create(u)

	if result.Error != nil {
		return 0, result.Error
	}
	return u.ID, nil
}

func (u *UserModel) Update() error {
	db := common.GetDB()

	return db.Save(u).Error
}
