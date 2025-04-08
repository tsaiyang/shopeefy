package di

import (
	"shopeefy/config"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/spf13/viper"
)

const algoshopKey = "algoshop"

type AppEnv struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectUri  string `mapstructure:"redirect_uri"`
	Scopes       string `mapstructure:"scopes"`
	ClientName   string `mapstructure:"client_name"`
	ClientHandle string `mapstructure:"client_handle"`
}

func InitShopifyAppEnv() *config.ShopifyApp {
	var env AppEnv
	if err := viper.UnmarshalKey(algoshopKey, &env); err != nil {
		panic(err)
	}

	app := config.ShopifyApp{
		App: goshopify.App{
			ApiKey:      env.ClientId,
			ApiSecret:   env.ClientSecret,
			Scope:       env.Scopes,
			RedirectUrl: env.RedirectUri,
		},
		ClientName:   env.ClientName,
		ClientHandle: env.ClientHandle,
	}

	return &app
}
