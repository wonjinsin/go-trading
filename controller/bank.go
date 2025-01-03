package controller

import (
	"context"
	"magmar/model/dto"
	"magmar/service"
	"magmar/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Bank ...
type Bank struct {
	BankSvc service.BankService
}

// NewBankController ...
func NewBankController(bankSvc service.BankService) BankController {
	return &Bank{
		BankSvc: bankSvc,
	}
}

// Deposit ...
func (b *Bank) Deposit(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow(util.LogCtrl)
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	var deposit dto.Deposit
	if err := c.Bind(&deposit); err != nil {
		zlog.With(intCtx).Warnw("BankSvc Deposit failed", "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	transaction, err := b.BankSvc.Deposit(intCtx, deposit)
	if err != nil {
		zlog.With(intCtx).Warnw("BankSvc Deposit failed", "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "Deposit OK", transaction)
}

// Withdrawal ...
func (b *Bank) Withdrawal(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow(util.LogCtrl)
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	var withdrawal dto.Withdrawal
	if err := c.Bind(&withdrawal); err != nil {
		zlog.With(intCtx).Warnw("BankSvc Withdrawal failed", "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	transaction, err := b.BankSvc.Withdrawal(intCtx, withdrawal)
	if err != nil {
		zlog.With(intCtx).Warnw("BankSvc Withdrawal failed", "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "Withdrawal OK", transaction)
}
