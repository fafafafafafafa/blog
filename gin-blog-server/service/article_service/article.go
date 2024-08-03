package article_service

import (
	"encoding/json"
	"go-gin-example/models"
	"go-gin-example/pkg/gredis"
	"go-gin-example/pkg/logging"
	"go-gin-example/service/cache_service"
)

type Article struct {
	ID            int
	TagId         int
	Title         string
	Desc          string
	Content       string
	CreatedBy     string
	ModifiedBy    string
	State         int
	CoverImageUrl string

	PageNum  int
	PageSize int
}

func (article *Article) GetArticlesTotal() (int, error) {
	maps := map[string]interface{}{
		"deleted_on": 0,
		"state":      article.State,
		"tag_id":     article.TagId,
	}
	return models.GetArticlesTotal(maps)
}
func (article *Article) GetArticles() ([]*models.Article, error) {
	// 首先查询缓存
	var cacheArticles []*models.Article
	cache_service_article := cache_service.Article{
		TagID:    article.TagId,
		State:    article.State,
		PageNum:  article.PageNum,
		PageSize: article.PageSize,
	}
	articleKey := cache_service_article.GetArticlesKey()
	if gredis.Exists(articleKey) {
		data, err := gredis.Get(articleKey)
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheArticles)
			if err != nil {
				return nil, err
			}
			return cacheArticles, nil
		}
	}
	// 查询数据库
	var articles []*models.Article
	var err error
	maps := map[string]interface{}{
		"deleted_on": 0,
		"state":      article.State,
		"tag_id":     article.TagId,
	}

	articles, err = models.GetArticles(article.PageNum, article.PageSize, maps)
	if err != nil {
		return nil, err
	}
	// 缓存
	if err := gredis.Set(articleKey, articles, 600); err != nil {
		logging.Warn(err)
	}
	return articles, nil
}

func (article *Article) ExistArticleByID() (bool, error) {

	return models.ExistArticleByID(article.ID)
}

func (article *Article) GetArticle() (*models.Article, error) {
	var cacheArticle models.Article
	cache_service_article := cache_service.Article{
		ID: article.ID,
	}
	key := cache_service_article.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
			return nil, err
		} else {
			err := json.Unmarshal(data, &cacheArticle)
			if err != nil {
				return nil, err
			}
			return &cacheArticle, nil
		}
	}
	var a *models.Article
	var err error
	a, err = models.GetArticle(map[string]interface{}{
		"deleted_on": 0,
		"id":         article.ID,
	})
	if err != nil {
		return nil, err
	}
	if err := gredis.Set(key, a, 600); err != nil {
		logging.Warn(err)
	}
	return a, nil
}

func (article *Article) AddArticle() error {
	maps := map[string]interface{}{
		"tag_id":          article.TagId,
		"title":           article.Title,
		"desc":            article.Desc,
		"content":         article.Content,
		"created_by":      article.CreatedBy,
		"state":           article.State,
		"cover_image_url": article.CoverImageUrl,
	}

	return models.AddArticle(maps)
}

func (article *Article) EditArticle() error {
	maps := map[string]interface{}{
		"tag_id":          article.TagId,
		"title":           article.Title,
		"desc":            article.Desc,
		"content":         article.Content,
		"modified_by":     article.ModifiedBy,
		"cover_image_url": article.CoverImageUrl,
	}
	if article.State >= 0 {
		maps["state"] = article.State
	}
	return models.EditArticle(article.ID, maps)
}

func (article *Article) DeleteArticle() error {
	return models.DeleteArticle(article.ID)
}
