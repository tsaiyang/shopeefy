package dao

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}

type gormUserDAO struct {
	db *gorm.DB
}

func (dao *gormUserDAO) Insert(ctx context.Context, user User) error {
	//TODO implement me
	panic("implement me")
}

func NewGormUserDAO(db *gorm.DB) UserDAO {
	return &gormUserDAO{db: db}
}

type User struct {
	Id   int64          `gorm:"primaryKey;autoIncrement"`
	Name sql.NullString `gorm:"type:varchar(20);uniqueIndex"`
}
