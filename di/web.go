package di

import (
	"shopeefy/internal/controller"

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

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func InitHandler(authHandler *controller.AuthHandler) []controller.Handler {
	return []controller.Handler{authHandler}
}
