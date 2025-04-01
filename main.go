package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/gin-gonic/gin"
	uuid "github.com/lithammer/shortuuid/v4"
)

var app = goshopify.App{
	ApiKey:      "964fbe18a1f4f7aa2239dddbbac693e5",
	ApiSecret:   "5a6be8f607a010611847610b3541f153",
	Scope:       "write_products",
	RedirectUrl: "https://7b6e-14-103-207-76.ngrok-free.app/auth/callback",
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")

	server.GET("/auth/auth", func(ctx *gin.Context) {
		shop := ctx.Query("shop")
		nonce := uuid.New()
		authUrl, err := app.AuthorizeUrl(shop, nonce)
		if err != nil {
			ctx.String(http.StatusOK, "fail to generate auth url")
			return
		}

		if ctx.Query("embedded") != "1" {
			ctx.Redirect(http.StatusFound, authUrl)
			return
		}

		if ctx.Query("escape") != "1" {
			redirectUri := fmt.Sprintf("https://%s%s?shop=%s&escape=1", ctx.Request.Host, ctx.Request.URL.Path, shop)
			ctx.HTML(http.StatusOK, "escape_iframe.html", gin.H{
				"ApiKey":      app.ApiKey,
				"Host":        ctx.Query("host"),
				"RedirectUri": redirectUri,
			})

			return
		}

		ctx.Redirect(http.StatusFound, authUrl)
	})

	server.GET("/auth/callback", func(ctx *gin.Context) {
		// shop := ctx.Query("shop")
		// code := ctx.Query("code")
		// token, err := app.GetAccessToken(ctx, shop, code)
		// if err != nil {
		// 	fmt.Println(err)
		// 	ctx.String(http.StatusOK, "fail to get access token")
		// 	return
		// }

		// ctx.JSON(http.StatusOK, gin.H{
		// 	"token": token,
		// 	"shop":  shop,
		// })
		client := createHTTPClient()
		resp, err := client.Get("https://sosososoko.myshopify.com")
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("请求超时:", err)
			} else {
				fmt.Println("请求失败:", err)
			}
		} else {
			defer resp.Body.Close()
			fmt.Println("请求成功，状态码:", resp.StatusCode)
		}
	})

	server.Run(":8080")
}
func createHTTPClient() *http.Client {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, // 添加 Shopify 推荐套件
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
		CurvePreferences:   []tls.CurveID{tls.X25519, tls.CurveP256}, // 添加曲线配置
		InsecureSkipVerify: false,                                    // 不建议在生产环境中使用
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 20 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   120 * time.Second,
	}
}
