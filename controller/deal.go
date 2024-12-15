package controller

import (
	"context"
	"magmar/service"
	"magmar/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Deal ...
type Deal struct {
	DealSvc service.DealService
}

// NewDealController ...
func NewDealController(dealSvc service.DealService) DealController {
	return &Deal{
		DealSvc: dealSvc,
	}
}

// Deal ...
func (d *Deal) Deal(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow(util.LogCtrl)
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	if err := d.DealSvc.Deal(intCtx); err != nil {
		zlog.With(intCtx).Warnw("DealSvc Deal failed", "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "Deal OK")
}
