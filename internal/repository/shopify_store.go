package repository

import (
	"context"
	"shopeefy/internal/model"
	"shopeefy/internal/repository/dao"
)

type ShopifyStoreRepo interface {
	FindStoreByDomain(ctx context.Context, domain string) (model.ShopifyStore, error)
	CreateStore(ctx context.Context, store model.ShopifyStore) error
}

type shopifyStoreRepo struct {
	shopifyStoreDAO dao.ShopifyStoreDAO
}

func (repo *shopifyStoreRepo) CreateStore(ctx context.Context, store model.ShopifyStore) error {
	return repo.shopifyStoreDAO.CreateStore(ctx, repo.toEntity(store))
}

func (repo *shopifyStoreRepo) FindStoreByDomain(ctx context.Context, domain string) (model.ShopifyStore, error) {
	store, err := repo.shopifyStoreDAO.FindStoreyDomain(ctx, domain)
	if err != nil {
		return model.ShopifyStore{}, err
	}

	return repo.toModel(store), nil
}

func NewShopifyStoreRepo(shopifyStoreDAO dao.ShopifyStoreDAO) ShopifyStoreRepo {
	return &shopifyStoreRepo{
		shopifyStoreDAO: shopifyStoreDAO,
	}
}

func (repo *shopifyStoreRepo) toModel(entity dao.ShopifyStore) model.ShopifyStore {
	return model.ShopifyStore{
		Id:          entity.Id,
		Domain:      entity.ShopifyStoreDomain,
		AccessToken: entity.AccessToken,
		Scopes:      entity.Scopes,
	}
}

func (repo *shopifyStoreRepo) toEntity(store model.ShopifyStore) dao.ShopifyStore {
	return dao.ShopifyStore{
		Id:                 store.Id,
		AppKey:             store.AppKey,
		ShopifyStoreDomain: store.Domain,
		AccessToken:        store.AccessToken,
		Scopes:             store.Scopes,
	}
}
