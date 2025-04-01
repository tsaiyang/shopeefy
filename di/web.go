package di

import (
	"shopeefy/internal/controller"

	"github.com/gin-gonic/gin"
)

func InitWebServer(middlewares []gin.HandlerFunc, handlers []controller.Handler) *gin.Engine {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")
	server.Use(middlewares...)

	for _, handler := range handlers {
		handler.RegisterRoutes(server)
	}

	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return nil
}

func InitHandler(userHandler *controller.UserHandler, authHandler *controller.AuthHandler) []controller.Handler {
	return []controller.Handler{userHandler, authHandler}
}
