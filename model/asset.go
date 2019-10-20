package model

import "github.com/jinzhu/gorm"

type Asset struct {
	gorm.Model
	UUID        string  `gorm:"type:varchar(255)" json:"uuid"`
	Name        string  `gorm:"type:varchar(255)" json:"name"`
	Description string  `json:"description"`
	CatalogID   int     `json:"catalogId"`
	Catalog     Catalog `json:"account"`
	PlacementId int     `json:"orderId"`
}
