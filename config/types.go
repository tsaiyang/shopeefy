package config

import goshopify "github.com/bold-commerce/go-shopify/v4"

type ShopifyApp struct {
	goshopify.App
	ClientName   string
	ClientHandle string
}
