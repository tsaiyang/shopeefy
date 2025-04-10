package controller

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(engine gin.IRouter)
}

type AnnouncementBarVO struct {
	Id         int64  `json:"id,omitempty"`
	Shop       string `json:"shop,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	ConfigAll  string `json:"config_all,omitempty"`
	ConfigPart string `json:"config_part,omitempty"`
	UpdateAt   int64  `json:"update_at,omitempty"`
}
