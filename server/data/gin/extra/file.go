package data

import (
	"github.com/flipped-aurora/gva/global"
	model "github.com/flipped-aurora/gva/model/gin/extra"
	"github.com/gookit/color"
	"gorm.io/gorm"
	"time"
)

var File = new(file)

type file struct{}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: exa_file_upload_and_downloads 表初始化数据
func (f *file) Init() error {
	files := []model.ExaFileUploadAndDownload{
		{global.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "10.png", "http://qmplusimg.henrongyi.top/gvalogo.png", "png", "158787308910.png"},
		{global.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "logo.png", "http://qmplusimg.henrongyi.top/1576554439myAvatar.png", "png", "1587973709logo.png"},
	}
	return global.Db.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 2}).Find(&[]model.ExaFileUploadAndDownload{}).RowsAffected == 2 {
			color.Danger.Println("\n[Mysql] --> exa_file_upload_and_downloads 表初始数据已存在!")
			return nil
		}
		if err := tx.Create(&files).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> exa_file_upload_and_downloads 表初始数据成功!")
		return nil
	})
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 定义表名
func (f *file) TableName() string {
	return "exa_file_upload_and_downloads"
}
