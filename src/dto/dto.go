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

type RegionGridDto struct {
	ID int
	RegionDescription string
}