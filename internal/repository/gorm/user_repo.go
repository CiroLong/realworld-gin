package gorm

import (
	"context"
	"errors"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"github/CiroLong/realworld-gin/internal/repository"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err == nil {
		return nil
	}

	// 唯一索引冲突（email / username）
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return repository.ErrUserAlreadyExist
	}

	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User

	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User

	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	res := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", user.ID).
		Updates(user)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return repository.ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) IsFollowing(ctx context.Context, followerID int64, followingID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepo) Follow(ctx context.Context, followerID int64, followingID int64) error {
	follow := &entity.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(follow).Error

	return err
}

func (r *UserRepo) UnFollow(ctx context.Context, followerID int64, followingID int64) error {
	err := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&entity.Follow{}).Error

	return err
}
