package service

import (
	"context"
	"github/CiroLong/realworld-gin/internal-v2/model/dto/user"
)

type UserService interface {
	// Register 用户注册
	Register(ctx context.Context, req *user.RegisterRequest) (*user.UserResponse, error)

	// Login 用户登录
	Login(ctx context.Context, req *user.LoginRequest) (*user.UserResponse, error)

	// GetCurrentUser 获取当前登录用户
	GetCurrentUser(ctx context.Context, userID int64) (*user.UserResponse, error)

	// UpdateCurrentUser 更新当前用户信息
	UpdateCurrentUser(ctx context.Context, userID int64, req *user.UpdateUserRequest) (*user.UserResponse, error)
}
