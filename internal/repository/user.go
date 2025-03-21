package repository

import (
	"shopeefy/internal/repository/cache"
	"shopeefy/internal/repository/dao"
)

type UserRepo interface {
}

type userRepo struct {
	userCache cache.UserCache
	userDAO   dao.UserDAO
}

var _ UserRepo = (*userRepo)(nil)

func NewUserRepo(userCache cache.UserCache, userDAO dao.UserDAO) UserRepo {
	return &userRepo{
		userCache: userCache,
		userDAO:   userDAO,
	}
}
