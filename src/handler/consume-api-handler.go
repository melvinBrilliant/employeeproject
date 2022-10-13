package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Photo struct {
	AlbumID int				`json:"albumId"`
	ID int					`json:"regionId"`
	Title string			`json:"title"`
	Url string				`json:"url"`
	ThumbnailUrl string		`json:"thumbnailUrl"`
}

func ConsumeApi(c echo.Context) error {
	response, err := http.Get("https://jsonplaceholder.typicode.com/photos")

	if (err != nil) {
		fmt.Println(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)

	if (err != nil) {
		fmt.Println(err.Error())
	}

	var photos []Photo
	json.Unmarshal(responseData, &photos)
	return c.JSON(http.StatusOK, photos)
}