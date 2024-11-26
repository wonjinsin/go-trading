package controller

import (
	"context"
	"magmar/model"
	"magmar/service"
	"magmar/util"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// User ...
type User struct {
	userSvc service.UserService
}

// NewUserController ...
func NewUserController(userSvc service.UserService) UserController {
	return &User{
		userSvc: userSvc,
	}
}

// GetUser ...
func (u *User) GetUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	uid := c.Param("uid")
	zlog.With(ctx).Infow("[New request]", "uid", uid)
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	if _, err = uuid.Parse(uid); err != nil {
		zlog.With(intCtx).Warnw("ID is not valid", "uid", uid, "err", err)
		return response(c, http.StatusBadRequest, "User is not valid")
	}

	user := &model.User{}
	if user, err = u.userSvc.GetUser(intCtx, uid); err != nil {
		zlog.With(intCtx).Warnw("UserSvc GetUser failed", "uid", uid, "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "GetUser OK", user)
}

// UpdateUser ...
func (u *User) UpdateUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow("[New request]")
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	uid := c.Param("uid")
	if _, err = uuid.Parse(uid); err != nil {
		zlog.With(intCtx).Warnw("ID is not valid", "uid", uid, "err", err)
		return response(c, http.StatusBadRequest, "User is not valid")
	}

	user := &model.User{}
	if err := c.Bind(user); err != nil {
		zlog.With(intCtx).Warnw("Bind error", "uid", uid, "user", user, "err", err)
		return response(c, http.StatusBadRequest, "Bind error")
	}
	if user, err = u.userSvc.UpdateUser(intCtx, uid, user); err != nil {
		zlog.With(intCtx).Errorw("UserSvc NewUser failed", "uid", uid, "user", user, "err", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "Update Deal OK", user)
}

// DeleteUser ...
func (u *User) DeleteUser(c echo.Context) (err error) {
	ctx := c.Request().Context()
	zlog.With(ctx).Infow("[New request]")
	intCtx, cancel := context.WithTimeout(ctx, util.CtxTimeOut)
	defer cancel()

	uid := c.Param("uid")
	if _, err = uuid.Parse(uid); err != nil {
		zlog.With(intCtx).Warnw("ID is not valid", "uid", uid, "err", err)
		return response(c, http.StatusBadRequest, "User is not valid")
	}

	if err = u.userSvc.DeleteUser(intCtx, uid); err != nil {
		zlog.With(intCtx).Errorw("UserSvc DeleteUser failed", "uid", uid, "err", err)
		return response(c, http.StatusInternalServerError, "DeleteUser failed")
	}

	return response(c, http.StatusOK, "Delete User OK")
}
