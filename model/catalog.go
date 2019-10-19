package model

import "github.com/jinzhu/gorm"

type Catalog struct {
	gorm.Model
	UUID        string `gorm:"type:varchar(255)" json:"uuid"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `json:"description"`
	AccountID   int
	Account     Account
}
