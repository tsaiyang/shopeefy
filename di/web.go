package di

import (
	"crypto/rand"
	"encoding/base64"
	"shopeefy/internal/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	return []gin.HandlerFunc{sessionMiddleware()}
}

func InitHandler(userHandler *controller.UserHandler, authHandler *controller.AuthHandler) []controller.Handler {
	return []controller.Handler{userHandler, authHandler}
}

func sessionMiddleware() gin.HandlerFunc {
	authKey, _ := generateSecretKey()
	encryptKey, _ := generateEncryptionKey()
	store := cookie.NewStore([]byte(authKey), []byte(encryptKey))

	return sessions.Sessions("shopify-sessid", store)
}

func generateSecretKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func generateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
