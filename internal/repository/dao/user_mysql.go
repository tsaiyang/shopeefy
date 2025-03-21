package dao

import (
	"context"
	"gorm.io/gorm"
)

type mysqlUserDAO struct {
	db *gorm.DB
}

func (dao *mysqlUserDAO) Insert(ctx context.Context, user User) error {
	//TODO implement me
	panic("implement me")
}

var _ UserDAO = (*mysqlUserDAO)(nil)

func NewGormUserDAO(db *gorm.DB) UserDAO {
	return &mysqlUserDAO{db: db}
}
