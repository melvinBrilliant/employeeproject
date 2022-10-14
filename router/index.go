package router

import (
	"com.melvin.employee/context"
	ctr "com.melvin.employee/controller"
	"github.com/labstack/echo"
	mdw "github.com/labstack/echo/middleware"
)

const base = "amfs-sales"

func Router() (e *echo.Echo) {
	e = echo.New()
	e.Use(mdw.CORS(), mdw.Logger())
	e.GET("/", ctr.SucessOnly)

	public := e.Group(base, context.MiddlewareChain)

	public.GET("/provincies", ctr.GetProvinsies)

	return e
}
