package data

import (
	"github.com/flipped-aurora/gva/global"
	model "github.com/flipped-aurora/gva/model/gf/system"
	"time"

	"github.com/gookit/color"

	"gorm.io/gorm"
)

var Authority = new(authority)

type authority struct{}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: authorities 表数据初始化
func (a *authority) Init() error {
	authorities := []model.Authority{
		{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "888", AuthorityName: global.I18n.T("OrdinaryUser"), ParentId: "0", DefaultRouter: "dashboard"},
		{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "8881", AuthorityName: global.I18n.T("NormalUserSubRole"), ParentId: "888", DefaultRouter: "dashboard"},
		{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "9528", AuthorityName: global.I18n.T("TestRole"), ParentId: "0", DefaultRouter: "dashboard"},
	}
	return global.Db.Transaction(func(tx *gorm.DB) error {
		if tx.Where("authority_id IN ? ", []string{"888", "9528"}).Find(&[]model.Authority{}).RowsAffected == 2 {
			color.Danger.Println("\n[Mysql] --> authorities 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&authorities).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	})
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 定义表名
func (a *authority) TableName() string {
	return "authorities"
}
