package handler

import (
	"net/http"
	// "strconv"

	"com.melvin.employee/src/entity"
	"github.com/labstack/echo/v4"
)

func GetAllRegions(c echo.Context) error {
	var region = new(entity.Region)
	regions := region.GetAll()
	return c.JSON(http.StatusOK, regions)
}

func FindRegionByIdQuery(c echo.Context) error {
	var region = new(entity.Region)
	id := c.QueryParam("id")
	if (id == "") {
		return GetAllRegions(c)
	}
	result := region.FindById(id)
	return c.JSON(http.StatusOK, result)
}

func FindRegionByIdPath(c echo.Context) error {
	id := c.Param("id")
	var region = new(entity.Region)
	result := region.FindById(id)
	return c.JSON(http.StatusOK, result)
}

func InsertRegion(c echo.Context) error {
	var regionInsert entity.Region
	var region = new(entity.Region)
	requestBody := c.Bind(&regionInsert)
	if (requestBody != nil) {
		return c.JSON(http.StatusUnprocessableEntity, "bad request")
	}
	return c.JSON(http.StatusAccepted, region.Save(regionInsert))
}

func UpdateRegion(c echo.Context) error {
	var regionUpdate entity.Region
	var region = new(entity.Region)
	requestBody := c.Bind(&regionUpdate)

	regionIsPresent := region.IsPresent(regionUpdate)
	if (!regionIsPresent) {
		return c.JSON(http.StatusNotFound, "Region not found")
	}

	if (requestBody == nil) {
		return c.JSON(http.StatusUnprocessableEntity, "bad request")
	}
	return c.JSON(http.StatusAccepted, regionUpdate.Save(regionUpdate))
}