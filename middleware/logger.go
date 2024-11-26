package middleware

import (
	"context"
	"magmar/util"
	"time"

	"github.com/labstack/echo/v4"
)

// SetTRID ...
func SetTRID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, util.TRID, util.GetTRID())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// RequestLogger ...
func RequestLogger(zlog *util.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()
			log := zlog.With(req.Context())

			fields := []interface{}{
				"status", res.Status,
				"latency", time.Since(start).String(),
				"method", req.Method,
				"uri", req.RequestURI,
				"host", req.Host,
				"remote_ip", c.RealIP(),
			}

			n := res.Status
			switch {
			case n >= 500:
				log.Infow("Server error", fields...)
			case n == 404 && req.RequestURI == "/favicon.ico":
				break
			case n >= 400:
				log.Infow("Client error", fields...)
			case n >= 300:
				log.Infow("Redirection", fields...)
			case n == 208:
				break
			default:
				log.Infow("Success", fields...)
			}

			return nil
		}
	}
}
