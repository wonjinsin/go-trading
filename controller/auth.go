package controller

import (
	"context"
	"magmar/model"
	"magmar/service"
	"magmar/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Auth ...
type Auth struct {
	authSvc service.AuthService
}

// NewAuthController ...
func NewAuthController(authSvc service.AuthService) AuthController {
	return &Auth{
		authSvc: authSvc,
	}
}

// Signup ...
func (a *Auth) Signup(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow("[New request]")
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	signup := &model.Signup{}
	if err := c.Bind(signup); err != nil {
		zlog.With(intCtx).Warnw("Bind error", "signup", signup, "err", err)
		return response(c, http.StatusBadRequest, err.Error())
	} else if !signup.Validate() {
		zlog.With(intCtx).Warnw("NewUser ValidateNewUser failed", "signup", signup)
		return response(c, http.StatusBadRequest, "Validate failed")
	}

	if user, err := a.authSvc.Signup(intCtx, signup); err != nil {
		zlog.With(intCtx).Errorw("AuthSvc NewUser failed", "user", user, "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "New Auth OK")
}

// Signin ...
func (a *Auth) Signin(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow("[New request]")
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	ctx = context.WithValue(ctx, util.LoginKey, true)
	defer cancel()

	signin := &model.Signin{}
	if err := c.Bind(signin); err != nil {
		zlog.With(intCtx).Warnw("Bind error", "signin", signin, "err", err)
		return response(c, http.StatusBadRequest, err.Error())
	} else if !signin.Validate() {
		zlog.With(intCtx).Warnw("Signin Validate failed", "signin", signin)
		return response(c, http.StatusBadRequest, "Validate failed")
	}

	token, err := a.authSvc.Signin(intCtx, signin)
	if err != nil {
		zlog.With(intCtx).Errorw("UserSvc Login failed", "signin", signin, "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "Signin OK", token)
}
