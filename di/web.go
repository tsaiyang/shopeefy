package di

import (
	"shopeefy/config"
	"shopeefy/internal/controller"
	"shopeefy/middlewares"

	"github.com/gin-gonic/gin"
)

func InitWebServer(middlewares []gin.HandlerFunc, handlers []controller.Handler) *gin.Engine {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")

	v1 := server.Group("/api/v1")
	v1.Use(middlewares...)

	for _, handler := range handlers {
		handler.RegisterRoutes(v1)
	}

	return server
}

func InitMiddlewares(app *config.ShopifyApp) []gin.HandlerFunc {
	jwtBuilder := middlewares.NewJwtSessionBuilder(app)
	return []gin.HandlerFunc{jwtBuilder.Build()}
}

func InitHandler(authHandler *controller.AuthHandler) []controller.Handler {
	return []controller.Handler{authHandler}
}
