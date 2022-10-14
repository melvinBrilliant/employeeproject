package controller

import "github.com/labstack/echo"

func GetProvinsies(c echo.Context) (err error) {
	var (
		ctx = newc(c)
	)

	return ctx.Success("Ok")
}
