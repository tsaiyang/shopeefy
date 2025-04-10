package service

import (
	"context"
	"shopeefy/internal/model"
	"shopeefy/internal/repository"
)

type AnnouncementBarService interface {
	PublishInDetail(ctx context.Context, bar model.AnnouncementBar) (int64, error)
	PublishInStatus(ctx context.Context, bid int64, shop string) error
	FindById(ctx context.Context, bid int64) (model.AnnouncementBar, error)
	Delete(ctx context.Context, bid int64, shop string) error
	Unpublish(ctx context.Context, bid int64, shop string) error
	GetByShop(ctx context.Context, shop string) ([]model.AnnouncementBar, error)
}

type announcementBarService struct {
	barRepo repository.AnnouncementBarRepo
}

func (service *announcementBarService) GetByShop(ctx context.Context, shop string) ([]model.AnnouncementBar, error) {
	return service.barRepo.GetByShop(ctx, shop)
}

func (service *announcementBarService) Unpublish(ctx context.Context, bid int64, shop string) error {
	return service.barRepo.UpdateStatus(ctx, bid, shop, model.AnnouncementBarStatusUnpublished)
}

func (service *announcementBarService) PublishInStatus(ctx context.Context, bid int64, shop string) error {
	return service.barRepo.UpdateStatus(ctx, bid, shop, model.AnnouncementBarStatusPublished)
}

func (service *announcementBarService) Delete(ctx context.Context, bid int64, shop string) error {
	return service.barRepo.UpdateStatus(ctx, bid, shop, model.AnnouncementBarStatusDeleted)
}

func (service *announcementBarService) FindById(ctx context.Context, bid int64) (model.AnnouncementBar, error) {
	return service.barRepo.FindById(ctx, bid)
}

func (service *announcementBarService) PublishInDetail(ctx context.Context, bar model.AnnouncementBar) (int64, error) {
	bar.Status = model.AnnouncementBarStatusPublished
	if bar.Id == 0 {
		return service.barRepo.Create(ctx, bar)
	}

	err := service.barRepo.Update(ctx, bar)
	return bar.Id, err
}

func NewAnnouncementBarService(barRepo repository.AnnouncementBarRepo) AnnouncementBarService {
	return &announcementBarService{
		barRepo: barRepo,
	}
}
