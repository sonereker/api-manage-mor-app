package model

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

func (a *Account) Disable() {
	a.Enabled = false
}

func (a *Account) Enable() {
	a.Enabled = true
}
