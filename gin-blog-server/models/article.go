package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model
	Tag Tag `json:"tag"`

	TagId         int    `json:"tag_id" gorm:"index"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
	CoverImageUrl string `json:"cover_image_url"`
}

// 获取文章列表
func GetArticles(pageNum int, pageSize int, maps map[string]interface{}) ([]*Article, error) {
	var articles []*Article
	// 查询articles时预加载tags

	if err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, nil
}

// 获取指定文章
func GetArticle(maps map[string]interface{}) (*Article, error) {
	var article Article
	if err := db.Where(maps).First(&article).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&article).Related(&article.Tag).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

// 新建文章
func AddArticle(maps map[string]interface{}) error {
	err := db.Create(&Article{

		TagId:         maps["tag_id"].(int),
		Title:         maps["title"].(string),
		Desc:          maps["desc"].(string),
		Content:       maps["content"].(string),
		CreatedBy:     maps["created_by"].(string),
		State:         maps["state"].(int),
		CoverImageUrl: maps["cover_image_url"].(string),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// 更新指定文章
func EditArticle(id int, data map[string]interface{}) error {
	err := db.Model(&Article{}).Where("id=?", id).Updates(data).Error
	return err
}

// 删除指定文章
func DeleteArticle(id int) error {
	err := db.Delete(&Article{}, id).Error
	return err
}

// // 回调函数
// func (article *Article) BeforeCreate(scope *gorm.Scope) {
// 	scope.SetColumn("created_on", time.Now().Unix())
// }

// // 回调函数
// func (article *Article) BeforeUpdate(scope *gorm.Scope) {
// 	scope.SetColumn("modified_on", time.Now().Unix())
// }

func GetArticlesTotal(maps map[string]interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func ExistArticleByID(id int) (bool, error) {

	var article Article
	if err := db.Select("id").Where("id=? AND deleted_on=?", id, 0).First(&article).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return article.ID > 0, nil
}

// 定时清理软删除的数据
func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}
