//go:build wireinject

package main

import (
	"shopeefy/di"
	"shopeefy/internal/controller"
	"shopeefy/internal/repository"
	"shopeefy/internal/repository/dao"
	"shopeefy/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third party
		di.InitShopifyAppEnv,
		di.InitDB,

		// dao
		dao.NewShopDAO,
		// cache

		// repository
		repository.NewShopRepo,

		// service
		service.NewShopService,

		// controller
		controller.NewAuthHandler,

		// server
		di.InitMiddlewares,
		di.InitHandler,
		di.InitWebServer,
	)

	return gin.Default()
}
