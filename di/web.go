package di

import (
	"github.com/gin-gonic/gin"
	"shopeefy/internal/controller"
)

func InitWebServer(middlewares []gin.HandlerFunc, handlers []controller.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)

	for _, handler := range handlers {
		handler.RegisterRoutes(server)
	}

	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return nil
}

func InitHandler(userHandler *controller.UserHandler) []controller.Handler {
	return []controller.Handler{userHandler}
}
