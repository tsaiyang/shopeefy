package service

import (
	"context"
	"shopeefy/internal/model"
	"shopeefy/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type ShopifyStoreService interface {
	GetAccessTokenByDomain(ctx context.Context, domain string) (string, error)
	CreateStore(ctx context.Context, store model.ShopifyStore) error
}

type shopifyStoreService struct {
	shopifyStoreRepo repository.ShopifyStoreRepo
}

func (service *shopifyStoreService) CreateStore(ctx context.Context, store model.ShopifyStore) error {
	encryptedAccessToken, err := bcrypt.GenerateFromPassword([]byte(store.AccessToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	store.AccessToken = string(encryptedAccessToken)
	return service.shopifyStoreRepo.CreateStore(ctx, store)
}

func (service *shopifyStoreService) GetAccessTokenByDomain(ctx context.Context, domain string) (string, error) {
	store, err := service.shopifyStoreRepo.FindStoreByDomain(ctx, domain)
	if err != nil {
		return "", err
	}

	return store.AccessToken, nil
}

func NewOAuthService(shopifyStoreRepo repository.ShopifyStoreRepo) ShopifyStoreService {
	return &shopifyStoreService{
		shopifyStoreRepo: shopifyStoreRepo,
	}
}
