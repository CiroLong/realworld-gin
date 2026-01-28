package service

import (
	"context"
	"errors"
	"github/CiroLong/realworld-gin/internal-v2/model/dto/user"
	"github/CiroLong/realworld-gin/internal-v2/model/entity"
	"github/CiroLong/realworld-gin/internal-v2/pkg/jwt"
	"github/CiroLong/realworld-gin/internal-v2/pkg/password"
	"github/CiroLong/realworld-gin/internal-v2/repository"
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

func (s *userService) Register(
	ctx context.Context,
	req *user.RegisterRequest,
) (*user.UserResponse, error) {

	// 1. email 是否已存在
	if _, err := s.userRepo.FindByEmail(ctx, req.User.Email); err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	// 2. username 是否已存在
	if _, err := s.userRepo.FindByUsername(ctx, req.User.Username); err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// 3. hash password
	hash, err := password.Hash(req.User.Password)
	if err != nil {
		return nil, err
	}

	// 4. 构造 Entity
	u := &entity.User{
		Email:        req.User.Email,
		Username:     req.User.Username,
		PasswordHash: hash,
		Bio:          "",
		Image:        "",
	}

	// 5. 持久化
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	// 6. 生成 JWT
	token, err := s.jwtMgr.Generate(u.ID)
	if err != nil {
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

func (s *userService) Login(
	ctx context.Context,
	req *user.LoginRequest,
) (*user.UserResponse, error) {

	// 1. 根据 email 查用户
	u, err := s.userRepo.FindByEmail(ctx, req.User.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// 2. 校验密码
	if !password.Verify(u.PasswordHash, req.User.Password) {
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

func (s *userService) GetCurrentUser(ctx context.Context, userID int64) (*user.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userService) UpdateCurrentUser(ctx context.Context, userID int64, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}
