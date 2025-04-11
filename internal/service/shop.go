package service

import (
	"context"
	"errors"
	"shopeefy/internal/model"
	"shopeefy/internal/repository"
)

var (
	ErrEmptyAccessToken = errors.New("access token is empty")
)

type ShopService interface {
	SaveAccessToken(ctx context.Context, shop model.Shop) error
	GetAccessTokenByShopName(ctx context.Context, shopName string) (string, error)
}

type shopService struct {
	shopRepo repository.ShopRepo
}

func (service *shopService) GetAccessTokenByShopName(ctx context.Context, shopName string) (string, error) {
	shop, err := service.shopRepo.FindByName(ctx, shopName)
	if err != nil {
		return "", err
	}
	if len(shop.AccessToken) == 0 {
		return "", ErrEmptyAccessToken
	}

	return shop.AccessToken, nil
}

func (service *shopService) SaveAccessToken(ctx context.Context, shop model.Shop) error {
	return service.shopRepo.Upsert(ctx, shop)
}

func NewShopService(shopRepo repository.ShopRepo) ShopService {
	return &shopService{
		shopRepo: shopRepo,
	}
}
