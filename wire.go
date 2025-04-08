//go:build wireinject

package main

import (
	"shopeefy/di"
	"shopeefy/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third party
		di.InitShopifyAppEnv,
		// dao

		// cache

		// repository

		// service

		// controller
		controller.NewAuthHandler,

		// server
		di.InitMiddlewares,
		di.InitHandler,
		di.InitWebServer,
	)

	return gin.Default()
}
