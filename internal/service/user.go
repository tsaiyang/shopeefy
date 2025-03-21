package service

import "shopeefy/internal/repository"

type UserService interface {
}

type userService struct {
	userRepo repository.UserRepo
}

var _ UserService = (*userService)(nil)

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}
