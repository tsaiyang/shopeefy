package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrBarAndShopUnmatched = errors.New("bar and shop unmatched")
)

type AnnouncementBarDAO interface {
	Insert(ctx context.Context, bar AnnouncementBar) (int64, error)
	UpdateById(ctx context.Context, bar AnnouncementBar) error
	GetById(ctx context.Context, bid int64) (AnnouncementBar, error)
	UpdateStatus(ctx context.Context, bid int64, shop string, status uint8) error
	GetByShop(ctx context.Context, shop string) ([]AnnouncementBar, error)
}

type gormAnnouncementBarDAO struct {
	db *gorm.DB
}

func (dao *gormAnnouncementBarDAO) GetByShop(ctx context.Context, shop string) (bars []AnnouncementBar, err error) {
	err = dao.db.WithContext(ctx).
		Select("id", "shop", "status", "config_part", "update_at").
		Where("shop = ? AND status in ?", shop, []uint8{1, 2}).
		Order("update_at DESC").
		Find(&bars).
		Error

	return
}

func (dao *gormAnnouncementBarDAO) UpdateStatus(ctx context.Context, bid int64, shop string, status uint8) error {
	res := dao.db.WithContext(ctx).
		Model(&AnnouncementBar{}).
		Where("id = ? AND shop = ?", bid, shop).
		Updates(map[string]any{
			"status":    status,
			"update_at": time.Now().Unix(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return ErrBarAndShopUnmatched
	}

	return nil
}

func (dao *gormAnnouncementBarDAO) UpdateById(ctx context.Context, bar AnnouncementBar) error {
	res := dao.db.WithContext(ctx).
		Model(&AnnouncementBar{}).
		Where("id = ? AND shop = ?", bar.Id).
		Updates(map[string]any{
			"status":      bar.Status,
			"config_all":  bar.ConfigAll,
			"config_part": bar.ConfigAll,
			"update_at":   time.Now().Unix(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrBarAndShopUnmatched
	}

	return nil
}

func (dao *gormAnnouncementBarDAO) GetById(ctx context.Context, bid int64) (bar AnnouncementBar, err error) {
	err = dao.db.WithContext(ctx).Where("id = ? AND status = 1", bid).First(&bar).Error
	return
}

func (dao *gormAnnouncementBarDAO) Insert(ctx context.Context, bar AnnouncementBar) (int64, error) {
	now := time.Now().Unix()
	bar.CreateAt = now
	bar.UpdateAt = now

	err := dao.db.WithContext(ctx).Create(&bar).Error
	return bar.Id, err
}

func NewAnnouncementBarDAO(db *gorm.DB) AnnouncementBarDAO {
	return &gormAnnouncementBarDAO{
		db: db,
	}
}

type AnnouncementBar struct {
	Id         int64  `gorm:"primaryKey;Increment"`
	Shop       string `gorm:"type:varchar(255);index;not null"`
	Status     uint8
	ConfigAll  string `gorm:"type:text;not null"`
	ConfigPart string `gorm:"type:varchar(255);not null"`
	UpdateAt   int64
	CreateAt   int64
}
