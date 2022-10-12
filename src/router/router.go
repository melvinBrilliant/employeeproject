package router // tempat mengumpulkan semua controller

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"com.melvin.employee/src/handler"
)

func New() *echo.Echo {
	e := echo.New()

	// set all middleware here
	e.Pre(middleware.Recover())
	// e.Pre(middleware.RemoveTrailingSlash())

	// set all routes here
	TerritoryRoutes(e)
	RegionRoutes(e)

	return e
}

func TerritoryRoutes(e *echo.Echo) {
	mapping := e.Group("/territory")
	mapping.GET("/", handler.GetAllTerritories)
	mapping.GET("/:id", handler.FindTerritoryByIdPath)
	mapping.GET("", handler.GetParam)

	mapping.POST("/", handler.InsertTerritory)
	mapping.PUT("/", handler.UpdateTerritory)
}

func RegionRoutes(e *echo.Echo) {
	mapping := e.Group("/region")
	mapping.GET("/", handler.GetAllRegions)
	mapping.GET("/:id", handler.FindRegionByIdPath)
	mapping.GET("", handler.FindRegionByIdQuery)

	mapping.POST("", handler.InsertRegion)
	mapping.PUT("", handler.UpdateRegion)
}

func EmployeeRoutes(e *echo.Group) {}

func RoleRoutes(e *echo.Group) {}