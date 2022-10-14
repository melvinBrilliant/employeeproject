package middleware

import (
	ct "com.melvin.employee/context"

	"github.com/labstack/echo"
)

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return ct.MiddlewareChain(next)(c)
	}
}
