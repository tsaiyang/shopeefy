package middlewares

import (
	"shopeefy/config"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtSessionBuilder struct {
	app          *config.ShopifyApp
	ignoredPaths []string
}

func NewJwtSessionBuilder(app *config.ShopifyApp) *JwtSessionBuilder {
	return &JwtSessionBuilder{
		app: app,
		// 对授权阶段的 url 不用验证
		ignoredPaths: []string{"/api/v1/auth/auth", "/api/v1/auth/callback"},
	}
}

func (builder *JwtSessionBuilder) AddIgnoredPath(path string) *JwtSessionBuilder {
	builder.ignoredPaths = append(builder.ignoredPaths, path)
	return builder
}

func (builder *JwtSessionBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(builder.ignoredPaths, ctx.Request.URL.Path) {
			return
		}

		// tokenStr := extraToken(ctx)
	}
}

func extraToken(ctx *gin.Context) string {
	bearer := ctx.GetHeader("Authorization")
	if len(bearer) == 0 {
		return ""
	}

	parts := strings.Split(bearer, " ")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

type SessionClaims struct {
	jwt.RegisteredClaims
	Iss  string `json:"iss"` // 签发者，通常是 "https://${shop-name}.myshopify.com/admin"
	Dest string `json:"dest"`
	Aud  string `json:"aud"` // 应用的API key
	Sub  string `json:"sub"`
	Exp  int64  `json:"exp"` // 过期时间（UNIX时间戳）
	Nbf  int64  `json:"nbf"` // 生效时间（UNIX时间戳）
	Iat  int64  `json:"iat"`
	Jti  string `json:"jti"`
	Sid  string `json:"sid"`
	Sig  string `json:"sig"`
}
