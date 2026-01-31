package repository

import (
	"context"
	"errors"
	"github/CiroLong/realworld-gin/internal/model/entity"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
)

// interface 接口与实现解耦
// 换 DB, 写 mock, 单测 Service
type UserRepo interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error

	// FindByEmail 根据 email 查找用户
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByUsername 根据 username 查找用户
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByID 根据 id 查找用户
	FindByID(ctx context.Context, id int64) (*entity.User, error)

	// Update 更新用户信息（部分字段）
	Update(ctx context.Context, user *entity.User) error

	IsFollowing(ctx context.Context, followerID int64, followingID int64) (bool, error)
	Follow(ctx context.Context, followerID int64, followingID int64) error
	UnFollow(ctx context.Context, followerID int64, followingID int64) error
}
