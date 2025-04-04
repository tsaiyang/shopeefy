// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"shopeefy/di"
	"shopeefy/internal/controller"
	"shopeefy/internal/repository"
	"shopeefy/internal/repository/cache"
	"shopeefy/internal/repository/dao"
	"shopeefy/internal/service"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	v := di.InitMiddlewares()
	cmdable := di.InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	db := di.InitDB()
	userDAO := dao.NewGormUserDAO(db)
	userRepo := repository.NewUserRepo(userCache, userDAO)
	userService := service.NewUserService(userRepo)
	userHandler := controller.NewUserHandler(userService)
	app := di.InitAlgoshopEnv()
	shopifyStoreDAO := dao.NewShopifyStoreDAO(db)
	shopifyStoreRepo := repository.NewShopifyStoreRepo(shopifyStoreDAO)
	shopifyStoreService := service.NewOAuthService(shopifyStoreRepo)
	authHandler := controller.NewAuthHandler(app, shopifyStoreService)
	v2 := di.InitHandler(userHandler, authHandler)
	engine := di.InitWebServer(v, v2)
	return engine
}
