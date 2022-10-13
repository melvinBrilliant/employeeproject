package entity

import (
	// "database/sql"

	"com.melvin.employee/src/dto"
)

func (Territory) TableName() string {
	return "Territory"
}

type Territory struct {
	ID string 						`gorm:"primaryKey;column:ID"`
	TerritoryDescription string 	`gorm:"column:TerritoryDescription"`
	RegionID int					`gorm:"column:RegionID"`
	Region Region 					`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreginkey:RegionID"`
}

func convertToTerritoryGridDto(territory Territory) dto.TerritoryGridDto {
	return dto.TerritoryGridDto{
		ID: territory.ID,
		TerritoryName: territory.TerritoryDescription,
		RegionID: territory.RegionID,
	}
}

func toTerritoryGridDtoList(territoryList []Territory) []dto.TerritoryGridDto {
	var result[] dto.TerritoryGridDto
	for i := 0; i < len(territoryList); i++ {
		result = append(result, convertToTerritoryGridDto(territoryList[i]))
	}
	return result
}

func toTerritoryRegionDto(territory Territory) dto.TerritoryRegionDto {
	return dto.TerritoryRegionDto {
		TerritoryID: territory.ID,
		RegionID: territory.RegionID,
		TerritoryName: territory.TerritoryDescription,
		RegionName: territory.Region.RegionDescription,
	}
}

func toTerritoryRegionDtoList(territories []Territory) []dto.TerritoryRegionDto {
	var result []dto.TerritoryRegionDto
	for i := 0; i < len(territories); i++ {
		result = append(result, toTerritoryRegionDto(territories[i]))
	}
	return result
}

func (Territory) GetAll() []dto.TerritoryGridDto {
	var territories []Territory
	db.Find(&territories) 
	return toTerritoryGridDtoList(territories)
}

func (Territory) GetAllEager() []Territory {
	var territories []Territory
	err := db.Model(&Territory{}).Preload("Region").Find(&territories).Error
	if (err != nil) {
		panic(err.Error())
	}
	return territories
}

func (Territory) GetAllUltraEager() []Territory {
	var territories []Territory
	err := db.Model(&Territory{}).Preload("Region.Territories").Find(&territories).Error
	if (err != nil) {
		panic(err.Error())
	}
	return territories
}

func (Territory) FindById(id string) dto.TerritoryGridDto {
	var territory Territory
	db.Where("\"ID\" = ?", id).First(&territory)
	return convertToTerritoryGridDto(territory)
}

func (Territory) FindByRegionId(id string) []dto.TerritoryRegionDto {
	var territories []Territory
	err := db.Model(&Territory{}).Preload("Region").Where("\"RegionID\" = ?", id).Find(&territories).Error
	if (err != nil) {
		panic(err.Error())		
	}
	return toTerritoryRegionDtoList(territories)
}

func (Territory) IsPresent(territory Territory) bool {
	isPresent := true
	if (&territory == &Territory{}) {
		isPresent = false
	}
	return isPresent
}

func (Territory) Save(territory Territory) string {
	db.Save(&territory)
	return "SUCCESS"
}

func (Territory) Delete(id string) {
	var territory = Territory{ID: id}
	db.Delete(&territory)
}