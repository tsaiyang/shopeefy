package controller

import (
	"net/http"
	"shopeefy/internal/model"
	"shopeefy/internal/service"
	"shopeefy/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnnouncementBarHandler struct {
	barService service.AnnouncementBarService
}

func NewAnnouncementBarHandler(barService service.AnnouncementBarService) *AnnouncementBarHandler {
	return &AnnouncementBarHandler{
		barService: barService,
	}
}

func (handler *AnnouncementBarHandler) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/a-bar")
	group.POST("/publish/detail", handler.PublishInDetail)
	group.POST("/publish/status", handler.PublishInStatus)
	group.POST("/delete", handler.Delete)
	group.POST("/unpublish", handler.Unpublish)
	group.POST("/list", handler.List)
	group.POST("/detail", handler.Detail)
}

func (handler *AnnouncementBarHandler) Detail(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	bar, err := handler.barService.FindById(ctx, req.Id)
	if err != nil {
		logger.Logger.Error("fail to find bar",
			zap.Int64("bar id", req.Id),
			zap.Error(err),
		)

		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	ctx.JSON(http.StatusOK, Result{Data: AnnouncementBarVO{
		Id:        bar.Id,
		Shop:      bar.Shop,
		ConfigAll: bar.ConfigAll,
		Status:    bar.Status.ToUint8(),
		UpdateAt:  bar.UpdateAt.Unix(),
	}})
}

func (handler *AnnouncementBarHandler) List(ctx *gin.Context) {
	var shop string

	bars, err := handler.barService.GetByShop(ctx, shop)
	if err != nil {
		logger.Logger.Error("fail to get bars by shop",
			zap.String("shop", shop),
			zap.Error(err),
		)

		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	res := make([]AnnouncementBarVO, 0, len(bars))
	for _, bar := range bars {
		res = append(res, AnnouncementBarVO{
			Id:         bar.Id,
			Shop:       bar.Shop,
			ConfigPart: bar.ConfigPart,
			Status:     bar.Status.ToUint8(),
			UpdateAt:   bar.UpdateAt.Unix(),
		})
	}

	ctx.JSON(http.StatusOK, Result{Data: res})
}

func (handler *AnnouncementBarHandler) Unpublish(ctx *gin.Context) {
	handler.changeStatus(ctx, barOpUnpublish)
}

func (handler *AnnouncementBarHandler) Delete(ctx *gin.Context) {
	handler.changeStatus(ctx, barOpDelete)
}

func (handler *AnnouncementBarHandler) PublishInStatus(ctx *gin.Context) {
	handler.changeStatus(ctx, barOpPublish)
}

func (handler *AnnouncementBarHandler) PublishInDetail(ctx *gin.Context) {
	var shop string

	type Req struct {
		Id         int64  `json:"id"`
		ConfigAll  string `json:"config_all"`
		ConfigPart string `json:"config_part"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	bid, err := handler.barService.PublishInDetail(ctx, model.AnnouncementBar{
		Id:         req.Id,
		Shop:       shop,
		ConfigAll:  req.ConfigAll,
		ConfigPart: req.ConfigPart,
	})
	if err != nil {
		logger.Logger.Error("save announcement bar failed",
			zap.String("shop", shop),
			zap.Error(err))

		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	ctx.JSON(http.StatusOK, Result{Data: bid})
}

func (handler *AnnouncementBarHandler) changeStatus(ctx *gin.Context, op barOp) {
	var shop string

	type Req struct {
		Id int64 `json:"id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	var err error
	if op == barOpDelete {
		err = handler.barService.Delete(ctx, req.Id, shop)
	} else if op == barOpUnpublish {
		err = handler.barService.Unpublish(ctx, req.Id, shop)
	} else if op == barOpPublish {
		err = handler.barService.PublishInStatus(ctx, req.Id, shop)
	}

	if err != nil {
		logger.Logger.Error("fail to "+string(op)+" announcement bar in status",
			zap.Int64("bar id", req.Id),
			zap.String("shop", shop),
			zap.Error(err),
		)

		ctx.JSON(http.StatusOK, Result{Code: serverErrCode, Msg: httpRespSystemError})
		return
	}

	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
}

type barOp string

const (
	barOpDelete    barOp = "delete"
	barOpPublish   barOp = "publish"
	barOpUnpublish barOp = "unpublish"
)
