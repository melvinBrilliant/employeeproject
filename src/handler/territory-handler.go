package handler

import (
	// "fmt"
	"net/http"
	// "strconv"

	"com.melvin.employee/src/dto"
	"com.melvin.employee/src/entity"
	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetAllTerritories(c echo.Context) error {
	var territory entity.Territory
	territories := territory.GetAll()
	return c.JSON(http.StatusOK, territories)
}

func GetAllTerritoriesEager(c echo.Context) error {
	var territory entity.Territory
	territories := territory.GetAllEager()
	return c.JSON(http.StatusOK, territories)
}

func GetAllTerritoriesUltraEager(c echo.Context) error {
	var territory entity.Territory
	territories := territory.GetAllUltraEager()
	return c.JSON(http.StatusOK, territories)
}

func FindWithQuery(c echo.Context) error {
	var territory entity.Territory
	id := c.QueryParam("id")
	regionId := c.QueryParam("region-id")

	if (id != "" && regionId != "") {
		return c.JSON(http.StatusBadRequest, "Bad Argument")
	} else if (id == "" && regionId == "") {
		return GetAllRegions(c)
	} else if (regionId == "") {
		return c.JSON(http.StatusOK, territory.FindById(id))
	} else if (id == "") {
		return c.JSON(http.StatusOK, territory.FindByRegionId(regionId))
	}

	result := territory.FindById(id)
	if (&result == &dto.TerritoryGridDto{}) {
		return c.JSON(http.StatusNotFound, "Territory not found")
	}

	return c.JSON(http.StatusOK, result)
}

func FindTerritoryByIdPath(c echo.Context) error {
	var territory entity.Territory
	id := c.Param("id")

	result := territory.FindById(id)
	if (&result == &dto.TerritoryGridDto{}) {
		return c.JSON(http.StatusNotFound, "Territory not found")
	}
	return c.JSON(http.StatusOK, result)
}

func FindTerritoryInRegion(c echo.Context) error {
	id := c.QueryParam("region-id")

	if id == "" {
		return GetAllTerritories(c)
	}

	var territory entity.Territory
	response := territory.FindByRegionId(id)

	return c.JSON(http.StatusOK, response)
}

func InsertTerritory(c echo.Context) error {
	var insertDto dto.TerritoryUpsertDto
	bindErr := c.Bind(&insertDto); if (bindErr != nil) {
		return c.JSON(http.StatusUnprocessableEntity, "bad request")
	}
	errValidation := validate.Struct(insertDto); if (errValidation != nil) {
		return c.JSON(http.StatusUnprocessableEntity, errValidation.Error())
	}
	var territory = entity.Territory{
		ID: insertDto.ID,
		TerritoryDescription: insertDto.TerritoryName,
		RegionID: insertDto.RegionID,
	}
	territory.Save(territory)
	return c.JSON(http.StatusAccepted, "New territory has been saved!")
}

func UpdateTerritory(c echo.Context) error {
	var updateDto dto.TerritoryUpsertDto
	bindErr := c.Bind(&updateDto); if (bindErr != nil) {
		return c.JSON(http.StatusUnprocessableEntity, "Bad request")
	}
	validateErr := validate.Struct(updateDto); if (validateErr != nil) {
		return c.JSON(http.StatusUnprocessableEntity, validateErr.Error())
	}
	var territory = entity.Territory{
		ID: updateDto.ID,
		TerritoryDescription: updateDto.TerritoryName,
		RegionID: updateDto.RegionID,
	}
	if (territory.IsPresent(territory)) {
		return c.JSON(http.StatusNotFound, "Territory not found")
	}
	// todo : region validation
	territory.Save(territory)
	return c.JSON(http.StatusAccepted, "Territory has been updated")
}

func DeleteTerritoryWithPath(c echo.Context) error {
	id := c.Param("id")
	var territory entity.Territory
	territory.Delete(id)

	return c.JSON(http.StatusAccepted, "Territory has been deleted")
}
