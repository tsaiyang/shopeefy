package controller

import (
	"fmt"
	"net/http"
	"shopeefy/config"
	"shopeefy/internal/model"
	"shopeefy/internal/service"
	"shopeefy/pkg/logger"
	"strings"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"go.uber.org/zap"
)

const (
	shopNameRegExp = `^[a-zA-Z0-9][a-zA-Z0-9\-]*\.myshopify\.com$`
)

const (
	clientErrCode = 400
	serverErrCode = 500

	httpRespInvalidShopName        = "invalid shop name"
	httpRespSpliceAuthUrlFail      = "fail to splice auth url"
	httpRespInvalidParams          = "invalid params"
	httpRespVerifyAuthUrlFailed    = "verify auth url failed"
	httpRespSystemError            = "system error"
	httpRespInvalidRequest         = "invalid request"
	httpRespFailToFetchAccessToken = "fail to fetch access token"
)

const (
	apiCallback =  "/api/v1/auth/callback"
)

type AuthHandler struct {
	app             *config.ShopifyApp
	stateCookieName string
	shopNameRegExp  *regexp.Regexp
	shopService     service.ShopService
}

func NewAuthHandler(app *config.ShopifyApp, shopService service.ShopService) *AuthHandler {
	return &AuthHandler{
		app:             app,
		stateCookieName: "jwt-state",
		shopNameRegExp:  regexp.MustCompile(shopNameRegExp, regexp.None),
		shopService:     shopService,
	}
}

func (handler *AuthHandler) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/auth")
	group.GET("/auth", handler.Auth2Url)
	group.GET("/callback", handler.Callback)
}

func (handler *AuthHandler) Auth2Url(ctx *gin.Context) {
	shop := ctx.Query("shop")
	matched, err := handler.shopNameRegExp.MatchString(shop)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}
	if !matched {
		ctx.JSON(http.StatusOK, Result{Code: clientErrCode, Msg: httpRespInvalidShopName})
		return
	}

	// 应该不存在或者过期了，但是目前拿到的 access_token 是 offline 模式，没有过期时间，所以暂时不校验是否过期
	// 未来 access_token 可能是 online 模式，到时候再说
	accessToken, err := handler.shopService.GetAccessTokenByShopName(ctx, shop)
	if err == nil && len(accessToken) > 0 {
		// 这里应该是返回到咱们的首页
		ctx.HTML(http.StatusFound, "index.html", gin.H{
			"Url": "https://www.baidu.com",
		})
		return
	}

	nonce := uuid.New()
	authUrl, err := handler.app.AuthorizeUrl(shop, nonce)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSpliceAuthUrlFail})
		return
	}

	// 非嵌入式的 app，直接重定向到授权页面，这个暂时不用看，因为咱现在的 app 都是嵌入式的
	if ctx.Query("embedded") != "1" {
		ok, err := handler.app.VerifyAuthorizationURL(ctx.Request.URL)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespVerifyAuthUrlFailed})
			return
		}
		if !ok {
			ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespInvalidParams})
			return
		}

		if err = handler.setStateCookie(ctx, nonce); err != nil {
			ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
			return
		}

		ctx.Redirect(http.StatusFound, authUrl)
		return
	}

	if ctx.Query("escape") != "1" {
		ok, err := handler.app.VerifyAuthorizationURL(ctx.Request.URL)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespVerifyAuthUrlFailed})
			return
		}
		if !ok {
			ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespInvalidParams})
			return
		}

		redirectUri := fmt.Sprintf("https://%s%s?shop=%s&escape=1&embedded=1", ctx.Request.Host, ctx.Request.URL.Path, shop)
		ctx.HTML(http.StatusOK, "escape_iframe.html", gin.H{
			"ApiKey":      handler.app.ApiKey,
			"Host":        ctx.Query("host"),
			"RedirectUri": redirectUri,
		})

		return
	}

	if err = handler.setStateCookie(ctx, nonce); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	fmt.Println(authUrl)
	ctx.Redirect(http.StatusFound, authUrl)
}

func (handler *AuthHandler) Callback(ctx *gin.Context) {
	validUrl, err := handler.app.VerifyAuthorizationURL(ctx.Request.URL)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespVerifyAuthUrlFailed})
		return
	}
	if !validUrl {
		logger.Logger.Error("invalid url", zap.Any("ctx.Request.URL", ctx.Request.URL))
		ctx.JSON(http.StatusOK, Result{Code: clientErrCode, Msg: httpRespInvalidParams})
		return
	}

	if err = handler.VerifyState(ctx); err != nil {
		logger.Logger.Error("fail to verify state", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: clientErrCode, Msg: httpRespInvalidRequest})
		return
	}

	shop := ctx.Query("shop")
	matched, err := handler.shopNameRegExp.MatchString(shop)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}
	if !matched {
		ctx.JSON(http.StatusOK, Result{Code: clientErrCode, Msg: httpRespInvalidShopName})
		return
	}

	code := ctx.Query("code")
	accessToken, err := handler.app.GetAccessToken(ctx, shop, code)
	if err != nil {
		logger.Logger.Error("fail to get access token", zap.String("shop: ", shop), zap.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespFailToFetchAccessToken})
		return
	}

	if err = handler.shopService.SaveAccessToken(ctx, model.Shop{
		Name:        shop,
		AccessToken: accessToken,
		IsActive:    true,
		Scope:       handler.app.Scope,
	}); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	shopName := strings.Split(shop, ".")[0]
	redirectUrl := fmt.Sprintf("https://admin.shopify.com/store/%s/apps/%s", shopName, handler.app.ClientHandle)

	ctx.Redirect(http.StatusFound, redirectUrl)
}

func (handler *AuthHandler) VerifyState(ctx *gin.Context) error {
	stateCookie, err := ctx.Cookie(handler.stateCookieName)
	fmt.Println(stateCookie)
	if err != nil {
		return fmt.Errorf("can't get state cookie %w", err)
	}

	var sc StateClaims
	_, err = jwt.ParseWithClaims(stateCookie, &sc, func(token *jwt.Token) (any, error) {
		return []byte(handler.app.ApiKey), nil
	})
	if err != nil {
		return fmt.Errorf("can't parse state cookie %w", err)
	}

	if sc.State != ctx.Query("state") {
		return fmt.Errorf("state mismatch")
	}

	return nil
}

func (handler *AuthHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{State: state}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(handler.app.ApiKey))
	if err != nil {
		return err
	}

	ctx.SetCookie(handler.stateCookieName, tokenString, 3600, apiCallback, "", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
