package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"shopeefy/config"
	"shopeefy/pkg/logger"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotAffected = errors.New("token not affected")
	ErrTokenInvalidAud  = errors.New("token invalid aud")
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

		tokenStr := extraToken(ctx)

		var sc SessionClaims
		parserOptions := []jwt.ParserOption{
			jwt.WithoutClaimsValidation(),
		}
		token, err := jwt.ParseWithClaims(tokenStr, &sc, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signature: %v", t.Header["alg"])
			}

			return []byte(builder.app.ApiSecret), nil
		}, parserOptions...)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err = builder.verifySessionClaims(&sc); err != nil {
			logger.Logger.Error("fail to verify session claims", zap.Error(err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("shop", sc.Dest)
	}
}

func (builder *JwtSessionBuilder) verifySessionClaims(claims *SessionClaims) error {
	now := time.Now().Unix()
	if claims.Exp < now {
		return ErrTokenExpired
	}
	if now < claims.Nbf {
		return ErrTokenNotAffected
	}
	if claims.Aud != builder.app.ApiKey {
		return ErrTokenInvalidAud
	}
	if !strings.Contains(claims.Iss, ".myshopify.com/admin") {
		return fmt.Errorf("无效的 iss 值: %s", claims.Iss)
	}

	return nil
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
