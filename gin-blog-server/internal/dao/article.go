package dao

import (
	"gin-blog/internal/model"

	"gorm.io/gorm"
)

const (
	STATUS_PUBLIC = iota + 1 // 公开
	STATUS_SECRET            // 私密
	STATUS_DRAFT             // 草稿
)

const (
	TYPE_ORIGINAL  = iota + 1 // 原创
	TYPE_REPRINT              // 转载
	TYPE_TRANSLATE            // 翻译
)

func GetArticleList(db *gorm.DB, page, size int, title string, isDelete *bool, status, typ, categoryId, tagId int) (list []model.Article, total int64, err error) {
	db = db.Model(model.Article{})

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if isDelete != nil {
		db = db.Where("is_delete", isDelete)
	}
	if status != 0 {
		db = db.Where("status", status)
	}
	if categoryId != 0 {
		db = db.Where("category_id", categoryId)
	}
	if typ != 0 {
		db = db.Where("type", typ)
	}

	db = db.Preload("Category").Preload("Tags").
		Joins("LEFT JOIN article_tag ON article_tag.article_id = article.id").
		Group("id") // 去重
	if tagId != 0 {
		db = db.Where("tag_id = ?", tagId)
	}

	result := db.Count(&total).
		Scopes(Paginate(page, size)).
		Order("is_top DESC, article.id DESC").
		Find(&list)
	return list, total, result.Error
}

// 新增/编辑文章, 同时根据 分类名称, 标签名称 维护关联表
func AddOrUpdateArticle(db *gorm.DB, article *model.Article, categoryName string, tagNames []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 分类不存在则创建
		category := model.Category{Name: categoryName}
		result := db.Model(&model.Category{}).Where("name", categoryName).FirstOrCreate(&category)
		if result.Error != nil {
			return result.Error
		}
		article.CategoryId = category.ID

		// 先 添加/更新 文章, 获取到其 ID
		if article.ID == 0 {
			result = db.Create(&article)
		} else {
			result = db.Model(&article).Where("id", article.ID).Updates(article)
		}
		if result.Error != nil {
			return result.Error
		}

		// 清空文章标签关联
		result = db.Delete(model.ArticleTag{}, "article_id", article.ID)
		if result.Error != nil {
			return result.Error
		}

		var articleTags []model.ArticleTag
		for _, tagName := range tagNames {
			// 标签不存在则创建
			tag := model.Tag{Name: tagName}
			result := db.Model(&model.Tag{}).Where("name", tagName).FirstOrCreate(&tag)
			if result.Error != nil {
				return result.Error
			}
			articleTags = append(articleTags, model.ArticleTag{
				ArticleId: article.ID,
				TagId:     tag.ID,
			})
		}
		result = db.Create(&articleTags)
		return result.Error
	})
}

func UpdateArticleTop(db *gorm.DB, id int, isTop bool) error {
	result := db.Model(&model.Article{Model: model.Model{ID: id}}).Update("is_top", isTop)
	return result.Error
}

func GetArticle(db *gorm.DB, id int) (data *model.Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(model.Article{Model: model.Model{ID: id}}).
		First(&data)
	return data, result.Error
}

// 软删除文章（修改）
func UpdateArticleSoftDelete(db *gorm.DB, ids []int, isDelete bool) (int64, error) {
	result := db.Model(model.Article{}).
		Where("id IN ?", ids).
		Update("is_delete", isDelete)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 物理删除文章
func DeleteArticle(db *gorm.DB, ids []int) (int64, error) {
	// 删除 [文章-标签] 关联
	result := db.Where("article_id IN ?", ids).Delete(&model.ArticleTag{})
	if result.Error != nil {
		return 0, result.Error
	}

	// 删除 [文章]
	result = db.Where("id IN ?", ids).Delete(&model.Article{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil

}

func ImportArticle(db *gorm.DB, userAuthId int, title, content, img string) error {
	article := model.Article{
		Title:   title,
		Content: content,
		Img:     img,
		Status:  STATUS_DRAFT,
		Type:    TYPE_ORIGINAL,
		UserId:  userAuthId,
	}

	result := db.Create(&article)
	return result.Error
}
