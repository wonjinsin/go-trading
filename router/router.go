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

	makeV1QaRoute(ver, svc)
}

func makeV1QaRoute(ver *echo.Group, svc *service.Service) {
	user := ver.Group("/qa")
	qaCt := ct.NewQaController(svc.Qa)
	user.POST("/ask", qaCt.Ask)
}
