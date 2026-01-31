package service

import (
	"context"
	"github/CiroLong/realworld-gin/internal/model/dto"
)

type UserService interface {
	// Register 用户注册
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)

	// Login 用户登录
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.UserResponse, error)

	// GetCurrentUser 获取当前登录用户
	GetCurrentUser(ctx context.Context, userID int64) (*dto.UserResponse, error)

	// UpdateCurrentUser 更新当前用户信息
	UpdateCurrentUser(ctx context.Context, userID int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error)

	FollowUserByName(ctx context.Context, userID int64, username string) (*dto.ProfileResponse, error)
	UnfollowUserByName(ctx context.Context, userID int64, username string) (*dto.ProfileResponse, error)

	// userID 传0 时不拿follow关系
	GetProfile(ctx context.Context, username string, userID int64) (*dto.ProfileResponse, error)
}
