package service

import (
	"context"
	"errors"
	"fmt"
	"github/CiroLong/realworld-gin/internal/model/dto"
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

func (s *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {

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
	return &dto.UserResponse{
		User: dto.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.UserResponse, error) {

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
	return &dto.UserResponse{
		User: dto.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

// authed
func (s *userService) GetCurrentUser(ctx context.Context, userID int64) (*dto.UserResponse, error) {
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
	return &dto.UserResponse{
		User: dto.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

// authed
func (s *userService) UpdateCurrentUser(ctx context.Context, userID int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
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
	return &dto.UserResponse{
		User: dto.UserDTO{
			Email:    u.Email,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    token,
		},
	}, nil
}

func (s *userService) FollowUserByName(ctx context.Context, userID int64, username string) (*dto.ProfileResponse, error) {
	// 1. 找目标用户
	target, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 不允许关注自己（可选，但很推荐）
	if target.ID == userID {
		return nil, errors.New("cannot follow yourself")
	}

	// 3. 建立 follow 关系（Repo 层保证幂等）
	if err := s.userRepo.Follow(ctx, userID, target.ID); err != nil {
		return nil, err
	}

	// 4. 返回 profile（following 一定是 true）
	return &dto.ProfileResponse{
		Profile: dto.ProfileDTO{
			Username:  target.Username,
			Bio:       target.Bio,
			Image:     target.Image,
			Following: true,
		},
	}, nil
}

func (s *userService) UnfollowUserByName(ctx context.Context, userID int64, username string) (*dto.ProfileResponse, error) {
	// 1. 找目标用户
	target, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 删除 follow 关系（幂等）
	if err := s.userRepo.UnFollow(ctx, userID, target.ID); err != nil {
		return nil, err
	}

	// 3. 返回 profile（following 一定是 false）
	return &dto.ProfileResponse{
		Profile: dto.ProfileDTO{
			Username:  target.Username,
			Bio:       target.Bio,
			Image:     target.Image,
			Following: false,
		},
	}, nil
}

// userID 传0时不需要读follow关系
func (s *userService) GetProfile(ctx context.Context, username string, userID int64) (*dto.ProfileResponse, error) {
	// 1. 查目标用户
	target, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 是否关注
	following := false
	if userID > 0 && userID != target.ID {
		following, _ = s.userRepo.IsFollowing(ctx, userID, target.ID)
	}

	// 3. 返回 profile
	return &dto.ProfileResponse{
		Profile: dto.ProfileDTO{
			Username:  target.Username,
			Bio:       target.Bio,
			Image:     target.Image,
			Following: following,
		},
	}, nil
}
