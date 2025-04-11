package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type AnnouncementBarActionDAO interface {
	Insert(ctx context.Context, action AnnouncementBarAction) error
}

type gormAnnouncementBarActionDAO struct {
	db *gorm.DB
}

func (dao *gormAnnouncementBarActionDAO) Insert(ctx context.Context, action AnnouncementBarAction) (err error) {
	now := time.Now().Unix()
	action.CreateAt = now
	action.UpdateAt = now

	err = dao.db.WithContext(ctx).Create(&action).Error
	return
}

func NewAnnouncementBarActionDAO(db *gorm.DB) AnnouncementBarActionDAO {
	return &gormAnnouncementBarActionDAO{
		db: db,
	}
}

type AnnouncementBarAction struct {
	Id    int64 `gorm:"primaryKey;autoIncrement"`
	Shop  string
	BarId int64
	// View-1 Click-2 Close-3
	ActionType uint8
	// 首页:1 商品详情页:2 collection页-3 购物车页面-4 其他-5
	PageType   uint8
	UserIp     string
	UserNation string
	UpdateAt   int64
	CreateAt   int64
}
