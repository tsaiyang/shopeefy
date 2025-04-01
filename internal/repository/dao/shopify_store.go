package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShopifyStoreDAO interface {
	FindStoreyDomain(ctx context.Context, domain string) (ShopifyStore, error)
	CreateStore(ctx context.Context, store ShopifyStore) error
}

type gormShopifyStoreDAO struct {
	db *gorm.DB
}

func (dao *gormShopifyStoreDAO) CreateStore(ctx context.Context, store ShopifyStore) error {
	now := time.Now().UnixMilli()
	store.CreateAt = now
	store.UpdateAt = now

	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"access_token":   store.AccessToken,
			"update_at":      now,
			"scopes":         store.Scopes,
			"uninstalled":    false,
			"uninstalled_at": 0,
		}),
	}).Create(&store).Error
}

func (dao *gormShopifyStoreDAO) FindStoreyDomain(ctx context.Context, domain string) (store ShopifyStore, err error) {
	err = dao.db.WithContext(ctx).Where("shopify_store_domain = ?", domain).First(&store).Error
	return
}

func NewShopifyStoreDAO(db *gorm.DB) ShopifyStoreDAO {
	return &gormShopifyStoreDAO{db: db}
}

// ShopifyStore 时间都用时间戳表示
type ShopifyStore struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`
	// 只取 app key 的前 5 位，应该就构成唯一了
	AppKey string `gorm:"type:TEXT;uniqueIndex:idx_app_domain,length:5"`
	// 商店的唯一域名（如 example.mystify.com）
	ShopifyStoreDomain string `gorm:"type=varchar(255);uniqueIndex:idx_app_domain,length:10;not null"`
	ShopifyUserId      int64
	// 加密存储的Token
	AccessToken string `gorm:"type:varchar(255)"`
	// 0 代表没过期时间
	ExpireAt      int64
	CreateAt      int64
	UpdateAt      int64
	Scopes        string `gorm:"type:text"`
	Uninstalled   bool
	UninstalledAt int64
	Nonce         string `gorm:"type:varchar(255)"`
}
