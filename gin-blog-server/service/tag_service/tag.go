package tag_service

import (
	"encoding/json"
	"go-gin-example/models"
	"go-gin-example/pkg/gredis"
	"go-gin-example/pkg/logging"
	"go-gin-example/service/cache_service"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (tag *Tag) ExistTagByID() (bool, error) {
	return models.ExistTagByID(tag.ID)
}

func (tag *Tag) GetTags() ([]*models.Tag, error) {
	// 查询缓存
	var cacheTags []*models.Tag
	cache_service_tag := cache_service.Tag{
		Name:     tag.Name,
		State:    tag.State,
		PageNum:  tag.PageNum,
		PageSize: tag.PageSize,
	}
	key := cache_service_tag.GetTagsKey()
	if gredis.Exists(key) {
		logging.Debug("从redis获取数据...")
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
			return nil, err
		} else {
			err := json.Unmarshal(data, &cacheTags)
			if err != nil {
				return nil, err
			}
			return cacheTags, nil
		}
	}
	// 查询数据库
	logging.Debug("从mysql获取数据...")
	maps := map[string]interface{}{
		"deleted_on": 0,
		"name":       tag.Name,
		"state":      tag.State,
	}
	var tags []*models.Tag
	var err error
	tags, err = models.GetTags(tag.PageNum, tag.PageSize, maps)
	if err != nil {
		return nil, err
	}

	if err := gredis.Set(key, tags, 600); err != nil {
		logging.Warn(err)
	}
	return tags, nil
}
func (tag *Tag) GetTagTotal() (int, error) {
	maps := map[string]interface{}{
		"deleted_on": 0,
		"name":       tag.Name,
		"state":      tag.State,
	}
	return models.GetTagTotal(maps)
}
func (tag *Tag) ExistTagByName() (bool, error) {
	return models.ExistTagByName(tag.Name)
}

func (tag *Tag) AddTag() error {
	return models.AddTag(tag.Name, tag.CreatedBy, tag.State)

}

func (tag *Tag) EditTag() error {
	maps := map[string]interface{}{
		"name":        tag.Name,
		"modified_by": tag.ModifiedBy,
		"deleted_on":  0,
	}
	if tag.State >= 0 {
		maps["state"] = tag.State
	}
	return models.EditTag(tag.ID, maps)

}

func (tag *Tag) DeleteTag() error {

	return models.DeleteTag(tag.ID)

}
