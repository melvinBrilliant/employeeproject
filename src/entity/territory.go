package entity

import (
	"database/sql"

	"com.melvin.employee/src/dto"
)

func (Territory) TableName() string {
	return "Territory"
}

type Territory struct {
	ID string 						`gorm:"primaryKey;column:ID"`
	TerritoryDescription string 	`gorm:"column:TerritoryDescription"`
	RegionID int 					`gorm:"column:RegionID"`
}

func (Territory) GetAll() []Territory {
	var territories []Territory
	db.Find(&territories)

	return territories
}

func (Territory) FindById(id string) Territory {
	var territory Territory
	db.Where("\"ID\" = ?", id).First(&territory)

	return territory
}

func (Territory) FindByRegionId(id string) []dto.TerritoryRegionDto {
	queryLine1 := "SELECT ter.\"ID\", reg.\"ID\", ter.\"TerritoryDescription\", reg.\"RegionDescription\" "
	queryLine2 := "FROM \"Territory\" AS ter "
	queryLine3 := "JOIN \"Region\" AS reg on reg.\"ID\" = ter.\"RegionID\" "
	queryLine4 := "WHERE ter.\"RegionID\" = @regionID "
	
	var dtos []dto.TerritoryRegionDto
	db.Raw(queryLine1 + queryLine2 + queryLine3 + queryLine4, sql.Named("regionID", id)).Scan(&dtos)
	return dtos
}

func (Territory) IsPresent(territory Territory) bool {
	isPresent := true
	if (territory == Territory{}) {
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