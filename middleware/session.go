package middleware

import (
	"context"
	"magmar/config"
	"magmar/util"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rbcervilla/redisstore/v8"
)

// SetSession ...
func SetSession(magmar *config.ViperConfig, zlog *util.Logger) (echo.MiddlewareFunc, error) {
	redisConn, err := util.RedisConnect(magmar, zlog)
	if err != nil {
		return nil, err
	}

	store, err := redisstore.NewRedisStore(context.Background(), redisConn)
	if err != nil {
		return nil, err
	}

	return session.Middleware(store), nil
}
