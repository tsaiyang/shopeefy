package controller

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(engine *gin.Engine)
}
