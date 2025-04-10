package middlewares

import "github.com/gin-gonic/gin"

type JwtSessionBuilder struct {
}

func NewJwtBuilder() *JwtSessionBuilder {
	return &JwtSessionBuilder{}
}

func (builder *JwtSessionBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	}
}
