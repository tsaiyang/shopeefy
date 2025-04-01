package controller

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"shopeefy/internal/model"
	"shopeefy/internal/service"
)

const (
	shopNameRegExp = `^[a-zA-Z0-9][a-zA-Z0-9\-]*\.myshopify\.com$`
)

const (
	clientErrCode = 4
	serverErrCode = 5

	httpRespInvalidShopName        = "invalid shop name"
	httpRespSpliceAuthUrlFail      = "fail to splice auth url"
	httpRespInvalidParams          = "invalid params"
	httpRespVerifyAuthUrlFailed    = "verify auth url failed"
	httpRespSystemError            = "system error"
	httpRespInvalidRequest         = "invalid request"
	httpRespFailToFetchAccessToken = "fail to fetch access token"
)

type AuthHandler struct {
	shopifyStoreService service.ShopifyStoreService
	app                 *goshopify.App
	stateCookieName     string
	shopNameRegExp      *regexp.Regexp
}

func NewAuthHandler(app *goshopify.App, shopifyStoreService service.ShopifyStoreService) *AuthHandler {
	return &AuthHandler{
		shopifyStoreService: shopifyStoreService,
		app:                 app,
		stateCookieName:     "jwt-state",
		shopNameRegExp:      regexp.MustCompile(shopNameRegExp, regexp.None),
	}
}

func (handler *AuthHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/auth")
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

	nonce := uuid.New()
	authUrl, err := handler.app.AuthorizeUrl(shop, nonce)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSpliceAuthUrlFail})
		return
	}

	// 非嵌入式的 app，直接重定向到授权页面
	if ctx.Query("embedded") != "1" {
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

		redirectUri := fmt.Sprintf("https://%s%s?shop=%s&escape=1", ctx.Request.Host, ctx.Request.URL.Path, shop)
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

	ctx.Redirect(http.StatusFound, authUrl)
}

func (handler *AuthHandler) Callback(ctx *gin.Context) {
	validUrl, err := handler.app.VerifyAuthorizationURL(ctx.Request.URL)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespVerifyAuthUrlFailed})
		return
	}
	if !validUrl {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespInvalidParams})
		return
	}

	if err := handler.VerifyState(ctx); err != nil {
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
	//accessToken, err := handler.app.GetAccessToken(ctx, shop, code)
	fmt.Println("shop: ", shop)
	fmt.Println("code: ", code)
	accessToken, err := handler.app.GetAccessToken(ctx, shop, code)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespFailToFetchAccessToken})
		return
	}

	// 保存 access token 到数据库
	if err = handler.shopifyStoreService.CreateStore(ctx, model.ShopifyStore{
		AccessToken: accessToken,
		AppKey:      handler.app.ApiKey,
		Domain:      shop,
		Scopes:      handler.app.Scope,
	}); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"shop":         shop,
		"access_token": accessToken,
	})
}

func (handler *AuthHandler) VerifyState(ctx *gin.Context) error {
	stateCookie, err := ctx.Cookie(handler.stateCookieName)
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

	ctx.SetCookie(handler.stateCookieName, tokenString, 3600, "/auth/callback", "", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

//func (handler *AuthHandler) exchangeToken(shop, code string) (string, error) {
//	client := &http.Client{}
//	reqURL := fmt.Sprintf("https://%s/admin/oauth/access_token", shop)
//
//	body := strings.NewReader(fmt.Sprintf(
//		"client_id=%s&client_secret=%s&code=%s",
//		handler.app.ApiKey,
//		handler.app.ApiSecret,
//		code,
//	))
//
//	req, _ := http.NewRequest("POST", reqURL, body)
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer func() { _ = resp.Body.Close() }()
//
//	if resp.StatusCode != http.StatusOK {
//		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	respBody, _ := io.ReadAll(resp.Body)
//
//	// 解析响应获取access_token（需要实现具体的JSON解析）
//	return parseAccessToken(respBody), nil
//}
//
//func parseAccessToken(data []byte) string {
//	type Token struct {
//		AccessToken string `json:"access_token"`
//		Scope       string `json:"scope"`
//	}
//
//	var token Token
//	if err := json.Unmarshal(data, &token); err != nil {
//		return ""
//	} else {
//		return token.AccessToken
//	}
//}
