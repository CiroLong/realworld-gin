package repository

import (
	"github/CiroLong/realworld-gin/internal/common"
	"github/CiroLong/realworld-gin/internal/models"
)

//	职责：
//	与数据库直接交互（CRUD 操作）
//	实现数据库查询优化（如分页查询）
//	使用 GORM/SQLx 等 ORM 工具
// 	关键点：
//	定义接口实现数据库解耦
//	使用 context.Context 传递超时控制

// 	提供分页、条件查询、事务等功能方法，不在上层拼 SQL。

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&models.UserModel{})
	db.AutoMigrate(&models.FollowModel{})
	db.AutoMigrate(&models.ArticleModel{})
	db.AutoMigrate(&models.TagModel{})
	db.AutoMigrate(&models.FavoriteModel{})
	db.AutoMigrate(&models.ArticleUserModel{})
	db.AutoMigrate(&models.CommentModel{})
}
