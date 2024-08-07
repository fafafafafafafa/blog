package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

func GetFrontStatistics(db *gorm.DB) (data model.FrontHomeVO, err error) {
	result := db.Model(&model.Article{}).Where("status = ? AND is_delete = ?", 1, 0).Count(&data.ArticleCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&model.UserAuth{}).Count(&data.UserCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&model.Message{}).Count(&data.MessageCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&model.Category{}).Count(&data.CategoryCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&model.Tag{}).Count(&data.TagCount)
	if result.Error != nil {
		return data, result.Error
	}

	data.Config, err = GetConfigMap(db)
	if err != nil {
		return data, err
	}

	return data, nil
}

// 前台文章列表（不在回收站并且状态为公开）
func GetBlogArticleList(db *gorm.DB, page, size, categoryId, tagId int) (data []model.Article, total int64, err error) {
	db = db.Model(model.Article{})
	db = db.Where("is_delete = 0 AND status = 1") // *

	if categoryId != 0 {
		db = db.Where("category_id = ?", categoryId)
	}
	if tagId != 0 {
		db = db.Where("id IN (SELECT article_id FROM article_tag WHERE tag_id = ?)", tagId)
	}

	db = db.Count(&total)
	result := db.Preload("Tags").Preload("Category").
		Order("is_top DESC, id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)

	return data, total, result.Error
}

// 前台文章详情（不在回收站并且状态为公开）
func GetBlogArticle(db *gorm.DB, id int) (data *model.Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(model.Article{Model: model.Model{ID: id}}).
		Where("is_delete = 0 AND status = 1"). // *
		First(&data)
	return data, result.Error
}

// 查询 n 篇推荐文章 (根据标签)
func GetRecommendList(db *gorm.DB, id, n int) (list []model.RecommendArticleVO, err error) {
	// sub1: 查出标签id列表
	// SELECT tag_id FROM `article_tag` WHERE `article_id` = ?
	sub1 := db.Table("article_tag").
		Select("tag_id").
		Where("article_id", id)
	// sub2: 查出这些标签对应的文章id列表 (去重, 且不包含当前文章)
	// SELECT DISTINCT article_id FROM (sub1) t
	// JOIN article_tag t1 ON t.tag_id = t1.tag_id
	// WHERE `article_id` != ?
	sub2 := db.Table("(?) t1", sub1).
		Select("DISTINCT article_id"). // 去重
		Joins("JOIN article_tag t ON t.tag_id = t1.tag_id").
		Where("article_id != ?", id)
	// 根据 文章id列表 查出文章信息 (前 n 个)
	result := db.Table("(?) t2", sub2).
		Select("id, title, img, created_at").
		Joins("JOIN article a ON t2.article_id = a.id").
		Where("a.is_delete = 0").
		Order("is_top, id DESC").
		Limit(n).
		Find(&list)
	return list, result.Error
}

// 查询最新的 n 篇文章
func GetNewestList(db *gorm.DB, n int) (data []model.RecommendArticleVO, err error) {
	result := db.Model(&model.Article{}).
		Select("id, title, img, created_at").
		Where("is_delete = 0 AND status = 1").
		Order("created_at DESC, id ASC").
		Limit(n).
		Find(&data)
	return data, result.Error
}

// 查询上一篇文章 (id < 当前文章 id)
func GetLastArticle(db *gorm.DB, id int) (val model.ArticlePaginationVO, err error) {
	sub := db.Table("article").Select("max(id)").Where("id < ?", id)
	result := db.Table("article").
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id = (?)", sub).
		Limit(1).
		Find(&val)
	return val, result.Error
}

// 查询下一篇文章 (id > 当前文章 id)
func GetNextArticle(db *gorm.DB, id int) (data model.ArticlePaginationVO, err error) {
	result := db.Model(&model.Article{}).
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id > ?", id).
		Limit(1).
		Find(&data)
	return data, result.Error
}

// 获取某篇文章的评论数
func GetArticleCommentCount(db *gorm.DB, articleId int) (count int64, err error) {
	result := db.Model(&model.Comment{}).
		Where("topic_id = ? AND type = 1 AND is_review = 1", articleId).
		Count(&count)
	return count, result.Error
}
