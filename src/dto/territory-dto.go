package dto

import ()

type TerritoryRegionDto struct {
	TerritoryID string `json:"territoryId"`
	RegionID int `json:"regionId"`
	TerritoryName string `json:"territoryName"`
	RegionName string `json:"regionName"`
}

type TerritoryGridDto struct {
	ID string
	TerritoryName string
	RegionID int
}

type TerritoryUpsertDto struct {
	ID string				`validate:"required"`
	TerritoryName string	`validate:"required"`
	RegionID int			`validate:"required"`
}