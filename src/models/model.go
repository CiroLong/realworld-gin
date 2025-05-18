package models

import "github/CiroLong/realworld-gin/src/common"

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&FollowModel{})
	db.AutoMigrate(&ArticleModel{})
	db.AutoMigrate(&TagModel{})
	db.AutoMigrate(&FavoriteModel{})
	db.AutoMigrate(&ArticleUserModel{})
	db.AutoMigrate(&CommentModel{})
}
