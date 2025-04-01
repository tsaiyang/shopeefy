package di

import (
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/spf13/viper"
)

const algoshopKey = "algoshop"

type AlgoshopEnv struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectUri  string `mapstructure:"redirect_uri"`
	Scopes       string `mapstructure:"scopes"`
}

func InitAlgoshopEnv() *goshopify.App {
	var env AlgoshopEnv
	if err := viper.UnmarshalKey(algoshopKey, &env); err != nil {
		panic(err)
	}

	app := goshopify.App{
		ApiKey:      env.ClientId,
		ApiSecret:   env.ClientSecret,
		RedirectUrl: env.RedirectUri,
		Scope:       env.Scopes,
	}

	return &app
}
