package controller

import "github.com/gin-gonic/gin"

type AuthHandler struct {
}

func (handler *AuthHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/auth")
	group.POST("/")
	group.Any("/callback")
}

func (handler *AuthHandler) Auth2Url(ctx *gin.Context) {

}

func (handler *AuthHandler) Callback(ctx *gin.Context) {

}
