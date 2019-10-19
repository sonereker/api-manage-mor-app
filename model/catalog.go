package model

import "github.com/jinzhu/gorm"

type Catalog struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)" json:"Name"`
	Description string `json:"Description"`
	AccountID   int
	Account     Account
}
