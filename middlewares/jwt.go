package middlewares

import "github.com/gin-gonic/gin"

type JwtBuilder struct {
}

func NewJwtBuilder() *JwtBuilder {
	return &JwtBuilder{}
}

func (builder *JwtBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	}
}
