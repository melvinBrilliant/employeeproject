package entity

import "com.melvin.employee/src/dto"

func (Region) TableName() string {
	return "Region"
}

type Region struct {
	ID int 						`gorm:"autoIncrement;column:ID"`
	RegionDescription string 	`gorm:"column:RegionDescription"`
	Territories []Territory		`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func toRegionGridDto(region Region) dto.RegionGridDto {
	return dto.RegionGridDto{
		ID: region.ID,
		RegionDescription: region.RegionDescription,
	}
}

func toRegionGridDtoList(regionList []Region) []dto.RegionGridDto {
	var result []dto.RegionGridDto
	for i := 0; i < len(regionList); i++ {
		result = append(result, toRegionGridDto(regionList[i]))
	}
	return result
}

func (Region) GetAll() []dto.RegionGridDto {
	var regions []Region
	db.Find(&regions)
	result := toRegionGridDtoList(regions)
	return result
}

func (Region) FindById(id string) dto.RegionGridDto {
	var region Region
	db.First(&region, id)
	return toRegionGridDto(region)
}

func (Region) FindAllEagerEntity() []Region {
	var regionList []Region
	err := db.Model(&Region{}).Preload("Territories").Find(&regionList).Error
	if (err != nil) {
		panic(err.Error())
	}
	return regionList
}

func (Region) IsPresent(region Region) bool {
	isPresent := true
	if (&region == &Region{}) {
		isPresent = false
	}
	return isPresent
}

func (Region) Save(region Region) string {
	db.Save(&region)
	return "SUCCESS"
}