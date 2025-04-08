package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShopDAO interface {
	Upsert(ctx context.Context, shop Shop) error
	FindByName(ctx context.Context, name string) (Shop, error)
}

type gormShopDAO struct {
	db *gorm.DB
}

func (dao *gormShopDAO) FindByName(ctx context.Context, name string) (shop Shop, err error) {
	err = dao.db.WithContext(ctx).Where("name = ?", name).First(&shop).Error
	return
}

func (dao *gormShopDAO) Upsert(ctx context.Context, shop Shop) error {
	now := time.Now().Unix()
	shop.UpdateAt = now
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]any{
			"access_token": shop.AccessToken,
			"scope":        shop.Scope,
			"expire_at":    shop.ExpireAt,
			"is_active":    shop.IsActive,
			"update_at":    now,
		}),
	}).Create(&shop).Error
}

func NewShopDAO(db *gorm.DB) ShopDAO {
	return &gormShopDAO{
		db: db,
	}
}

type Shop struct {
	Id          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	AccessToken string `gorm:"type:varchar(255)"`
	IsActive    bool   `gorm:"type:bool;default:true"`
	Scope       string `gorm:"type:varchar(1000)"`
	ExpireAt    int64
	UpdateAt    int64
	CreateAt    int64
}
