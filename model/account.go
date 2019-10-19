package model

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `json:"Password"`
	Enabled  bool   `json:"Enabled"`
}

func (a *Account) Disable() {
	a.Enabled = false
}

func (a *Account) Enable() {
	a.Enabled = true
}
