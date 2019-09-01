package model

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Email      string `gorm:"unique" json:"email"`
	Password   string `json:"password"`
	CustomerID int    `json:"customerId"`
	Enabled    bool   `json:"statusId"`
}

func (a *Account) Disable() {
	a.Enabled = false
}

func (a *Account) Enable() {
	a.Enabled = true
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Account{})
	return db
}
