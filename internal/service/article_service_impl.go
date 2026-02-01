package service

import (
	"context"
	"errors"
	"github/CiroLong/realworld-gin/internal/model/dto"
	"github/CiroLong/realworld-gin/internal/model/entity"
	"github/CiroLong/realworld-gin/internal/pkg/common"
	"github/CiroLong/realworld-gin/internal/pkg/utils"
	"github/CiroLong/realworld-gin/internal/repository"
)

type articleService struct {
	articleRepo repository.ArticleRepo
	userRepo    repository.UserRepo
}

func NewArticleService(articleRepo repository.ArticleRepo, userRepo repository.UserRepo) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (s articleService) CreateArticle(ctx context.Context, authorID int64, req *dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	// 1. 生成 Article 实体
	articleEntity := &entity.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    authorID,
	}
	// 2. 调用 Repo 创建文章（DB 操作全部 Repo 完成）
	if err := s.articleRepo.Create(ctx, articleEntity); err != nil {
		return nil, err
	}

	// 3. 处理 tag（Repo 提供接口）
	var tags []*entity.Tag
	if len(req.Article.TagList) > 0 {
		var err error
		tags, err = s.articleRepo.GetOrCreateTags(ctx, req.Article.TagList)
		if err != nil {
			return nil, err
		}

		if err = s.articleRepo.ReplaceArticleTags(ctx, articleEntity.ID, tags); err != nil {
			return nil, err
		}
	}

	// 4. 获取作者信息
	author, err := s.userRepo.FindByID(ctx, authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	// 5. 拼装 DTO（业务层逻辑）
	authorDTO := dto.AuthorDTO{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: false, // 用户不可以关注自己
	}

	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	// 创建时没有收藏
	return dto.NewArticleResponse(articleEntity, tagNames, authorDTO, false), nil
}

// GetArticle 获取单篇文章
func (s articleService) GetArticle(ctx context.Context, slug string) (*dto.ArticleResponse, error) {
	// 1. 获取文章
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, common.ErrNotFound
	}

	// 2. 获取作者
	author, err := s.userRepo.FindByID(ctx, article.AuthorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	authorDTO := dto.AuthorDTO{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: false, // TODO：GetArticle不一定带Auth
	}

	// 3. 获取标签
	tags, err := s.articleRepo.GetTagsByArticleID(ctx, article.ID)
	if err != nil {
		return nil, err
	}
	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	// 4. 获取收藏信息
	favorited := false
	// TODO: 如果这里传入了userID应该进行查询

	// 5. 返回 DTO
	return dto.NewArticleResponse(article, tagNames, authorDTO, favorited), nil
}

// UpdateArticle 更新
func (s articleService) UpdateArticle(ctx context.Context, slug string, userID int64, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	// 1. 查article
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, common.ErrNotFound
	}

	// 无权限更新
	if article.AuthorID != userID {
		return nil, common.ErrPermissionDenied
	}

	// 只更新非空字段
	if req.Article.Title != "" {
		article.Title = req.Article.Title
		article.Slug = utils.GenerateSlug(req.Article.Title)
	}
	if req.Article.Description != "" {
		article.Description = req.Article.Description
	}
	if req.Article.Body != "" {
		article.Body = req.Article.Body
	}

	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}

	// 获取标签和作者信息
	author, _ := s.userRepo.FindByID(ctx, article.AuthorID)
	authorDTO := dto.AuthorDTO{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: false, // 不可以关注自己
	}

	tags, _ := s.articleRepo.GetTagsByArticleID(ctx, article.ID)
	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	// 收藏信息
	favorited, err := s.articleRepo.IsFavorited(ctx, userID, article.ID)
	if err != nil {
		return nil, err
	}

	return dto.NewArticleResponse(article, tagNames, authorDTO, favorited), nil
}

// DeleteArticle 删除
func (s articleService) DeleteArticle(ctx context.Context, slug string, userID int64) error {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return common.ErrNotFound
	}

	if article.AuthorID != userID {
		return common.ErrPermissionDenied
	}

	return s.articleRepo.Delete(ctx, article.ID)
}

// FavoriteArticle 点赞
func (s articleService) FavoriteArticle(ctx context.Context, slug string, userID int64) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, common.ErrNotFound
	}

	favorited, err := s.articleRepo.IsFavorited(ctx, userID, article.ID)
	if err != nil {
		return nil, err
	}
	if !favorited {
		// 添加收藏
		if err := s.articleRepo.AddFavorite(ctx, userID, article.ID); err != nil {
			return nil, err
		}
	}
	// 获取最新 favoritesCount
	article.FavoritesCount, _ = s.articleRepo.CountFavorites(ctx, article.ID)

	// 获取作者和标签
	author, _ := s.userRepo.FindByID(ctx, article.AuthorID)
	// 获取follow关系
	follow, err := s.userRepo.IsFollowing(ctx, userID, article.AuthorID)
	if err != nil {
		return nil, err
	}

	authorDTO := dto.AuthorDTO{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: follow,
	}

	tags, _ := s.articleRepo.GetTagsByArticleID(ctx, article.ID)
	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	return dto.NewArticleResponse(article, tagNames, authorDTO, true), nil
}

// UnfavoriteArticle 取消点赞
func (s articleService) UnfavoriteArticle(ctx context.Context, slug string, userID int64) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, common.ErrNotFound
	}
	favorited, err := s.articleRepo.IsFavorited(ctx, userID, article.ID)
	if err != nil {
		return nil, err
	}

	if favorited {
		// 移除收藏
		if err = s.articleRepo.RemoveFavorite(ctx, userID, article.ID); err != nil {
			return nil, err
		}
	}

	// 获取最新 favoritesCount
	article.FavoritesCount, _ = s.articleRepo.CountFavorites(ctx, article.ID)

	// 获取作者和标签
	author, _ := s.userRepo.FindByID(ctx, article.AuthorID)
	// 获取follow关系
	follow, err := s.userRepo.IsFollowing(ctx, userID, article.AuthorID)
	if err != nil {
		return nil, err
	}
	authorDTO := dto.AuthorDTO{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: follow,
	}

	tags, _ := s.articleRepo.GetTagsByArticleID(ctx, article.ID)
	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	return dto.NewArticleResponse(article, tagNames, authorDTO, false), nil
}

func (s articleService) ListArticles(
	ctx context.Context,
	tag string,
	author string,
	favorited string,
	userID int64,
	limit int,
	offset int,
) (*dto.MultipleArticlesResponse, error) {
	// 1. 构造过滤条件
	filter := repository.ListArticlesFilter{
		Limit:  limit,
		Offset: offset,
	}
	if tag != "" {
		filter.Tag = &tag
	}
	if author != "" {
		filter.Author = &author
	}
	if favorited != "" {
		filter.FavoritedBy = &favorited
	}

	// 2. 查文章列表 + 总数
	articles, total, err := s.articleRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return &dto.MultipleArticlesResponse{
			Articles:      make([]dto.ArticleWithoutBodyDTO, 0),
			ArticlesCount: 0,
		}, nil
	}

	// 3. 批量准备 articleID
	articleIDs := make([]int64, 0, len(articles))
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
	}

	// 4. 批量获取 tag
	tagsMap, err := s.articleRepo.GetTagsByArticleIDs(ctx, articleIDs)
	if err != nil {
		return nil, err
	}

	// 5. 拼 DTO
	articleDTOs := make([]dto.ArticleWithoutBodyDTO, 0, len(articles))

	for _, a := range articles {
		// 作者
		authorEntity, err := s.userRepo.FindByID(ctx, a.AuthorID)
		if err != nil {
			return nil, err
		}

		following := false
		favoritedFlag := false

		if userID > 0 {
			following, _ = s.userRepo.IsFollowing(ctx, userID, a.AuthorID)
			favoritedFlag, _ = s.articleRepo.IsFavorited(ctx, userID, a.ID)
		}

		articleDTOs = append(articleDTOs, dto.ArticleWithoutBodyDTO{
			Slug:           a.Slug,
			Title:          a.Title,
			Description:    a.Description,
			TagList:        tagsMap[a.ID],
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
			FavoritesCount: a.FavoritesCount,
			Favorited:      favoritedFlag,
			Author: dto.AuthorDTO{
				Username:  authorEntity.Username,
				Bio:       authorEntity.Bio,
				Image:     authorEntity.Image,
				Following: following,
			},
		})
	}

	return &dto.MultipleArticlesResponse{
		Articles:      articleDTOs,
		ArticlesCount: int(total),
	}, nil
}

func (s articleService) FeedArticles(
	ctx context.Context,
	userID int64,
	limit int,
	offset int,
) (*dto.MultipleArticlesResponse, error) {
	articles, total, err := s.articleRepo.Feed(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return &dto.MultipleArticlesResponse{
			Articles:      make([]dto.ArticleWithoutBodyDTO, 0),
			ArticlesCount: 0,
		}, nil
	}

	articleIDs := make([]int64, 0, len(articles))
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
	}

	tagsMap, err := s.articleRepo.GetTagsByArticleIDs(ctx, articleIDs)
	if err != nil {
		return nil, err
	}

	articleDTOs := make([]dto.ArticleWithoutBodyDTO, 0, len(articles))

	for _, a := range articles {
		author, err := s.userRepo.FindByID(ctx, a.AuthorID)
		if err != nil {
			return nil, err
		}
		favoritedFlag := false

		if userID > 0 {
			favoritedFlag, _ = s.articleRepo.IsFavorited(ctx, userID, a.ID)
		}

		articleDTOs = append(articleDTOs, dto.ArticleWithoutBodyDTO{
			Slug:           a.Slug,
			Title:          a.Title,
			Description:    a.Description,
			TagList:        tagsMap[a.ID],
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
			Favorited:      favoritedFlag,
			FavoritesCount: a.FavoritesCount,
			Author: dto.AuthorDTO{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: true,
			},
		})
	}

	return &dto.MultipleArticlesResponse{
		Articles:      articleDTOs,
		ArticlesCount: int(total),
	}, nil
}

func (s articleService) ListTags(ctx context.Context) ([]string, error) {
	return s.articleRepo.ListTags(ctx)
}
