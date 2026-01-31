package entity

import "time"

type Follow struct {
	FollowerID  int64 `gorm:"primaryKey"`
	FollowingID int64 `gorm:"primaryKey"`

	CreatedAt time.Time
}
