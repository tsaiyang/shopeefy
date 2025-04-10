package repository

import (
	"context"
	"shopeefy/internal/model"
	"shopeefy/internal/repository/dao"
	"time"
)

type AnnouncementBarRepo interface {
	Create(ctx context.Context, bar model.AnnouncementBar) (int64, error)
	Update(ctx context.Context, bar model.AnnouncementBar) error
	UpdateStatus(ctx context.Context, bid int64, shop string, status model.AnnouncementBarStatus) error
	FindById(ctx context.Context, bid int64) (model.AnnouncementBar, error)
	GetByShop(ctx context.Context, shop string) ([]model.AnnouncementBar, error)
}

type announcementBarRepo struct {
	barDAO dao.AnnouncementBarDAO
}

func (repo *announcementBarRepo) GetByShop(ctx context.Context, shop string) ([]model.AnnouncementBar, error) {
	bars, err := repo.barDAO.GetByShop(ctx, shop)
	if err != nil {
		return nil, err
	}

	res := make([]model.AnnouncementBar, 0, len(bars))
	for _, bar := range bars {
		res = append(res, repo.toModel(bar))
	}

	return res, nil
}

func (repo *announcementBarRepo) UpdateStatus(ctx context.Context, bid int64, shop string, status model.AnnouncementBarStatus) error {
	return repo.barDAO.UpdateStatus(ctx, bid, shop, uint8(status))
}

func (repo *announcementBarRepo) Update(ctx context.Context, bar model.AnnouncementBar) error {
	return repo.barDAO.UpdateById(ctx, repo.toEntity(bar))
}

func (repo *announcementBarRepo) FindById(ctx context.Context, bid int64) (model.AnnouncementBar, error) {
	bar, err := repo.barDAO.GetById(ctx, bid)
	if err != nil {
		return model.AnnouncementBar{}, err
	}

	return repo.toModel(bar), nil
}

func (repo *announcementBarRepo) Create(ctx context.Context, bar model.AnnouncementBar) (int64, error) {
	return repo.barDAO.Insert(ctx, repo.toEntity(bar))
}

func NewAnnouncementBarRepo(barDAO dao.AnnouncementBarDAO) AnnouncementBarRepo {
	return &announcementBarRepo{
		barDAO: barDAO,
	}
}

func (repo *announcementBarRepo) toEntity(bar model.AnnouncementBar) dao.AnnouncementBar {
	return dao.AnnouncementBar{
		Id:         bar.Id,
		Shop:       bar.Shop,
		ConfigAll:  bar.ConfigAll,
		ConfigPart: bar.ConfigPart,
		Status:     bar.Status.ToUint8(),
		CreateAt:   bar.CreateAt.Unix(),
		UpdateAt:   bar.UpdateAt.Unix(),
	}
}

func (repo *announcementBarRepo) toModel(bar dao.AnnouncementBar) model.AnnouncementBar {
	return model.AnnouncementBar{
		Id:         bar.Id,
		Shop:       bar.Shop,
		ConfigAll:  bar.ConfigAll,
		ConfigPart: bar.ConfigPart,
		Status:     model.AnnouncementBarStatus(bar.Status),
		UpdateAt:   time.Unix(bar.UpdateAt, 0),
		CreateAt:   time.Unix(bar.CreateAt, 0),
	}
}
