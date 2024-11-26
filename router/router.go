package router

import (
	ct "magmar/controller"
	"magmar/service"

	"github.com/labstack/echo/v4"
)

// Init ...
func Init(e *echo.Echo, svc *service.Service) {
	api := e.Group("/api")
	ver := api.Group("/v1")

	makeV1AuthRoute(ver, svc)
	makeV1UserRoute(ver, svc)
}

func makeV1Route(ver *echo.Group, svc *service.Service) {
	makeV1UserRoute(ver, svc)
}

func makeV1AuthRoute(ver *echo.Group, svc *service.Service) {
	user := ver.Group("/auth")
	authCt := ct.NewAuthController(svc.Auth)
	user.POST("/signup", authCt.Signup)
	user.POST("/signin", authCt.Signin)
}

func makeV1UserRoute(ver *echo.Group, svc *service.Service) {
	user := ver.Group("/user")
	userCt := ct.NewUserController(svc.User)
	user.GET("/:uid", userCt.GetUser)
	user.PUT("/:uid", userCt.UpdateUser)
	user.DELETE("/:uid", userCt.DeleteUser)
}
