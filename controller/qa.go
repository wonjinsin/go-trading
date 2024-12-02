package controller

import (
	"context"
	"magmar/service"
	"magmar/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Qa ...
type Qa struct {
	QaSvc service.QaService
}

// NewQaController ...
func NewQaController(qaSvc service.QaService) QaController {
	return &Qa{
		QaSvc: qaSvc,
	}
}

// Ask ...
func (q *Qa) Ask(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow("[New request]")
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	if err := q.QaSvc.Ask(intCtx); err != nil {
		zlog.With(intCtx).Errorw("QaSvc Ask failed", "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "Ask OK")
}
