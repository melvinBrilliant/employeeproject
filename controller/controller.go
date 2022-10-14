package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

const (
	ApplicationJson   string = "application/json"
	MultipartFormData string = "multipart/form-data"
	FormUrlEncoded    string = "application/x-www-form-urlencoded"
	ContentType       string = "Content-Type"
	ContentMaxMemory  int64  = 1024 << 10
)

type ContextIndex uint

const (
	SessionContext ContextIndex = 10000
)

type EchoResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func SucessOnly(c echo.Context) (err error) {
	return c.JSON(http.StatusOK,
		EchoResponse{
			Status: "berhasil",
			Data:   nil,
		})
}
