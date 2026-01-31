package service

import (
	"context"
	"errors"
	"fmt"
	"github/CiroLong/realworld-gin/internal/model/dto/user"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"github/CiroLong/realworld-gin/internal/pkg/jwt"
	"github/CiroLong/realworld-gin/internal/pkg/password"
	"github/CiroLong/realworld-gin/internal/repository"
	"log"
)

type userService struct {
	userRepo repository.UserRepo
	jwtMgr   jwt.Manager
}

func NewUserService(
	userRepo repository.UserRepo,
	jwtMgr jwt.Manager,
) UserService {
	return &userService{
		userRepo: userRepo,
		jwtMgr:   jwtMgr,
	}
}

func (s *userService) Register(ctx context.Context, req *user.RegisterRequest) (*user.UserResponse, error) {

	// 1. email 是否已存在
	if _, err := s.userRepo.FindByEmail(ctx, req.User.Email); err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		log.Println("err", err.Error())
		return nil, err
	}

	// 2. username 是否已存在
	if _, err := s.userRepo.FindByUsername(ctx, req.User.Username); err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		log.Println("err", err.Error())
		return nil, err
	}

	// 3. hash password
	hash, err := password.Hash(req.User.Password)
	if err != nil {
		return nil, err
	}

	// 4. 构造 Entity
	u := &entity.User{
		Email:    req.User.Email,
		Username: req.User.Username,
		Password: hash,
		Bio:      "",
		Image:    "",
	}

	// 5. 持久化
	if err := s.userRepo.Create(ctx, u); err != nil {
		fmt.Errorf("create fail")
		return nil, err
	}

	// 6. 生成 JWT
	token, err := s.jwtMgr.Generate(u.ID)
	if err != nil {
		fmt.Errorf("Generate jwt fail")
		return nil, err
	}

	// 7. 组装 Response DTO
	return &user.UserResponse{
		User: user.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

func (s *userService) Login(ctx context.Context, req *user.LoginRequest) (*user.UserResponse, error) {

	// 1. 根据 email 查用户
	u, err := s.userRepo.FindByEmail(ctx, req.User.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// 2. 校验密码
	if !password.Verify(u.Password, req.User.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 3. 生成 JWT
	token, err := s.jwtMgr.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	// 4. 返回 Response DTO
	return &user.UserResponse{
		User: user.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

// authed
func (s *userService) GetCurrentUser(ctx context.Context, userID int64) (*user.UserResponse, error) {
	// 1. 根据 userID 查用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 生成 token（RealWorld 要求返回 token）
	token, err := s.jwtMgr.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	// 3. 组装响应 DTO
	return &user.UserResponse{
		User: user.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

// authed
func (s *userService) UpdateCurrentUser(ctx context.Context, userID int64, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	// 1. 查当前用户（确保存在）
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 2. 按需更新字段
	if req.User.Email != nil {
		u.Email = *req.User.Email
	}

	if req.User.Username != nil {
		u.Username = *req.User.Username
	}

	if req.User.Bio != nil {
		u.Bio = *req.User.Bio
	}

	if req.User.Image != nil {
		u.Image = *req.User.Image
	}

	if req.User.Password != nil {
		hashed, err := password.Hash(*req.User.Password)
		if err != nil {
			return nil, err
		}
		u.Password = hashed
	}

	// 3. 更新数据库
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 4. 生成 token（可轮换）
	token, err := s.jwtMgr.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	// 5. 响应
	return &user.UserResponse{
		User: user.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}
