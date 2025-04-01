//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"shopeefy/di"
	"shopeefy/internal/controller"
	"shopeefy/internal/repository"
	"shopeefy/internal/repository/cache"
	"shopeefy/internal/repository/dao"
	"shopeefy/internal/service"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third party
		di.InitDB,
		di.InitRedis,
		di.InitAlgoshopEnv,

		// dao
		dao.NewGormUserDAO,
		dao.NewShopifyStoreDAO,

		// cache
		cache.NewRedisUserCache,

		// repository
		repository.NewUserRepo,
		repository.NewShopifyStoreRepo,

		// service
		service.NewUserService,
		service.NewOAuthService,

		// controller
		controller.NewUserHandler,
		controller.NewAuthHandler,

		// server
		di.InitMiddlewares,
		di.InitHandler,
		di.InitWebServer,
	)

	return gin.Default()
}
