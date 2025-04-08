package repository

import (
	"context"
	"shopeefy/internal/model"
	"shopeefy/internal/repository/dao"
	"time"
)

type ShopRepo interface {
	Upsert(ctx context.Context, shop model.Shop) error
	FindByName(ctx context.Context, name string) (model.Shop, error)
}

type shopRepo struct {
	shopDAO dao.ShopDAO
}

func (repo *shopRepo) FindByName(ctx context.Context, name string) (model.Shop, error) {
	shop, err := repo.shopDAO.FindByName(ctx, name)
	if err != nil {
		return model.Shop{}, err
	}

	return repo.toModel(shop), nil
}

func (repo *shopRepo) Upsert(ctx context.Context, shop model.Shop) error {
	return repo.shopDAO.Upsert(ctx, repo.toEntity(shop))
}

func NewShopRepo(shopDAO dao.ShopDAO) ShopRepo {
	return &shopRepo{
		shopDAO: shopDAO,
	}
}

func (repo *shopRepo) toModel(shop dao.Shop) model.Shop {
	return model.Shop{
		Id:          shop.Id,
		Name:        shop.Name,
		AccessToken: shop.AccessToken,
		IsActive:    shop.IsActive,
		Scope:       shop.Scope,
		ExpireAt:    time.Unix(shop.ExpireAt, 0),
		UpdateAt:    time.Unix(shop.UpdateAt, 0),
		CreateAt:    time.Unix(shop.CreateAt, 0),
	}
}

func (repo *shopRepo) toEntity(shop model.Shop) dao.Shop {
	return dao.Shop{
		Id:          shop.Id,
		Name:        shop.Name,
		AccessToken: shop.AccessToken,
		IsActive:    shop.IsActive,
		Scope:       shop.Scope,
		ExpireAt:    shop.ExpireAt.Unix(),
		UpdateAt:    shop.UpdateAt.Unix(),
		CreateAt:    shop.CreateAt.Unix(),
	}
}
