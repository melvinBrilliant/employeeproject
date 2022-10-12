package entity

import ()

func (Region) TableName() string {
	return "Region"
}

type Region struct {
	ID int 						`gorm:"autoIncrement;column:ID"`
	RegionDescription string 	`gorm:"column:RegionDescription"`
}

func (Region) GetAll() []Region {
	var regions []Region
	db.Find(&regions)

	return regions
}

func (Region) FindById(id string) Region {
	var region Region
	db.First(&region, id)
	return region
}

func (Region) IsPresent(region Region) bool {
	isPresent := true
	if (region == Region{}) {
		isPresent = false
	}
	return isPresent
}

func (Region) Save(region Region) string {
	db.Save(&region)
	return "SUCCESS"
}