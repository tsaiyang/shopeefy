package controller

import (
	"github.com/gin-gonic/gin"
	"shopeefy/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
}
