package gorm

import (
	"context"
	"errors"
	"fmt"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"github/CiroLong/realworld-gin/internal/pkg/common"
	"github/CiroLong/realworld-gin/internal/pkg/utils"
	"github/CiroLong/realworld-gin/internal/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) repository.ArticleRepo {
	return &articleRepo{db: db}
}

func (a articleRepo) FindBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	var article entity.Article
	err := a.db.WithContext(ctx).Where("slug = ?", slug).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	return &article, err
}

func (a articleRepo) Create(ctx context.Context, article *entity.Article) error {
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	article.Slug = utils.GenerateSlug(article.Title)
	const maxRetry = 3

	for i := 0; i < maxRetry; i++ {
		// 尝试创建
		err := a.db.WithContext(ctx).Create(article).Error
		if err == nil {
			return nil
		}
		// 检查是否是 slug 冲突
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "slug") {
			// 生成新的 slug 并重试
			article.Slug = utils.GenerateSlug(article.Title)
			continue
		}
		// 其他错误直接返回
		return err
	}
	return fmt.Errorf("failed to create article after %d retries due to slug conflict", maxRetry)

}

func (a articleRepo) Update(ctx context.Context, article *entity.Article) error {
	article.UpdatedAt = time.Now()
	return a.db.WithContext(ctx).Save(article).Error
}

func (a articleRepo) Delete(ctx context.Context, articleID int64) error {
	return a.db.WithContext(ctx).Delete(&entity.Article{}, articleID).Error
}

func (a articleRepo) GetTagsByArticleID(ctx context.Context, articleID int64) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := a.db.WithContext(ctx).Model(&entity.ArticleTag{}).
		Select("tags.id, tags.name").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Where("article_tags.article_id = ?", articleID).
		Scan(&tags).Error
	return tags, err
}

// GetOrCreateTags 获取或创建具有给定名称的标签列表，确保每个标签的唯一性。
func (a articleRepo) GetOrCreateTags(ctx context.Context, names []string) ([]*entity.Tag, error) {
	if len(names) == 0 {
		return []*entity.Tag{}, nil
	}

	tags := make([]*entity.Tag, 0, len(names))
	for _, name := range names {
		var tag entity.Tag
		// 使用 FirstOrCreate 保证唯一性
		err := a.db.WithContext(ctx).Where("name = ?", name).FirstOrCreate(&tag, entity.Tag{Name: name}).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

// ReplaceArticleTags 替换文章的标签，使用事务确保删除旧标签和添加新标签操作的原子性。
func (a articleRepo) ReplaceArticleTags(ctx context.Context, articleID int64, tags []*entity.Tag) error {
	// 使用gorm的事务操作保证原子性
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1️. 删除旧的关联
		if err := tx.Where("article_id = ?", articleID).Delete(&entity.ArticleTag{}).Error; err != nil {
			return err
		}

		// 2. 插入新的关联
		if len(tags) == 0 {
			return nil
		}
		articleTags := make([]*entity.ArticleTag, len(tags))
		for i, tag := range tags {
			articleTags[i] = &entity.ArticleTag{ArticleID: articleID, TagID: tag.ID}
		}
		return tx.Create(&articleTags).Error
	})
}

func (a articleRepo) IsFavorited(ctx context.Context, userID, articleID int64) (bool, error) {
	var count int64
	err := a.db.WithContext(ctx).
		Model(&entity.Favorite{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a articleRepo) AddFavorite(ctx context.Context, userID, articleID int64) error {
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fav entity.Favorite
		err := tx.Where("article_id = ? AND user_id = ?", articleID, userID).First(&fav).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在才插入
			if err := tx.Create(&entity.Favorite{ArticleID: articleID, UserID: userID}).Error; err != nil {
				return err
			}
			// 增加计数
			return tx.Model(&entity.Article{}).Where("id = ?", articleID).Update("favorites_count", gorm.Expr("favorites_count + ?", 1)).Error
		} else if err != nil {
			return err
		}
		// 已存在，不操作计数
		return nil
	})
}

func (a articleRepo) RemoveFavorite(ctx context.Context, userID, articleID int64) error {
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&entity.Favorite{}).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Article{}).
			Where("id = ? AND favorites_count > 0", articleID).
			Update("favorites_count", gorm.Expr("favorites_count - ?", 1)).Error
	})
}

func (a articleRepo) CountFavorites(ctx context.Context, articleID int64) (int, error) {
	var article entity.Article
	if err := a.db.WithContext(ctx).Select("favorites_count").Where("id = ?", articleID).First(&article).Error; err != nil {
		return 0, err
	}
	return article.FavoritesCount, nil
}

func (a articleRepo) List(ctx context.Context, query repository.ListArticlesFilter) ([]*entity.Article, int64, error) {
	db := a.db.WithContext(ctx).Model(&entity.Article{})

	// --- tag 过滤 ---
	if query.Tag != nil {
		db = db.Joins(
			"JOIN article_tags at ON at.article_id = articles.id").
			Joins(
				"JOIN tags t ON t.id = at.tag_id").
			Where("t.name = ?", *query.Tag)
	}

	// --- author 过滤 ---
	if query.Author != nil {
		db = db.Joins(
			"JOIN users u ON u.id = articles.author_id").
			Where("u.username = ?", *query.Author)
	}

	// --- favorited 过滤 ---
	if query.FavoritedBy != nil {
		db = db.Joins(
			"JOIN favorites f ON f.article_id = articles.id").
			Joins(
				"JOIN users u2 ON u2.id = f.user_id").
			Where("u2.username = ?", *query.FavoritedBy)
	}

	// --- 统计总数 ---
	var total int64
	if err := db.
		Distinct("articles.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// --- 查询文章 ---
	var articles []*entity.Article
	if err := db.
		Select("articles.*").
		Order("articles.created_at DESC").
		Limit(query.Limit).
		Offset(query.Offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (a articleRepo) Feed(ctx context.Context, userID int64, limit, offset int) ([]*entity.Article, int64, error) {
	db := a.db.WithContext(ctx).
		Model(&entity.Article{}).
		Joins(
			"JOIN follows f ON f.following_id = articles.author_id").
		Where("f.follower_id = ?", userID)

	// 总数
	var total int64
	if err := db.
		Distinct("articles.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []*entity.Article
	if err := db.
		Order("articles.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (a articleRepo) ListTags(ctx context.Context) ([]string, error) {
	var tags []string

	if err := a.db.WithContext(ctx).
		Model(&entity.Tag{}).
		Distinct("name").
		Order("name ASC").
		Pluck("name", &tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (a articleRepo) GetTagsByArticleIDs(ctx context.Context, articleIDs []int64) (map[int64][]string, error) {
	type row struct {
		ArticleID int64
		TagName   string
	}

	var rows []row

	if err := a.db.WithContext(ctx).
		Table("article_tags").
		Select("article_tags.article_id, tags.name AS tag_name").
		Joins("JOIN tags ON tags.id = article_tags.tag_id").
		Where("article_tags.article_id IN ?", articleIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[int64][]string, len(articleIDs))
	for _, r := range rows {
		result[r.ArticleID] = append(result[r.ArticleID], r.TagName)
	}

	return result, nil
}
