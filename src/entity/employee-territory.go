package entity

import ()

func (EmployeeTerritory) TableName() string {
	return "EmployeeTerritory"
}

type EmployeeTerritory struct {
	EmployeeID *Employee `gorm:"primaryKey"`
	TerritoryID *Territory `gorm:"primaryKey"`
}

func (EmployeeTerritory) GetAll() []EmployeeTerritory {
	var employeeTerritories []EmployeeTerritory
	db.Find(&employeeTerritories)

	return employeeTerritories
}