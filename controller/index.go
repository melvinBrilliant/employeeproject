package controller

import (
	ct "com.melvin.employee/context"

	echo "github.com/labstack/echo"
)

const emptyString = ""

func newc(c echo.Context) *ct.Context {
	return ct.NewContext(c)
}
