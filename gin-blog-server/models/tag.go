package models

import "github.com/jinzhu/gorm"

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]*Tag, error) {
	var tags []*Tag
	if err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func AddTag(name string, created_by string, state int) error {
	var tag Tag = Tag{
		Name:      name,
		CreatedBy: created_by,
		State:     state,
	}

	err := db.Create(&tag).Error
	return err
}

func EditTag(id int, data map[string]interface{}) error {
	err := db.Model(&Tag{}).Where("id=?", id).Updates(data).Error
	return err
}

func DeleteTag(id int) error {
	err := db.Delete(&Tag{}, id).Error
	return err
}

// // 回调函数
// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("CreatedOn", time.Now().Unix())
// 	return nil
// }

// // 回调函数
// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	return nil
// }

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	var err error
	if err = db.Select("id").Where("id=? AND deleted_on=?", id, 0).First(&tag).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return tag.ID > 0, nil
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	if err := db.Select("id").Where("name=? AND deleted_on=?", name, 0).First(&tag).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return tag.ID > 0, nil
}

// 定时清理软删除的数据
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}
