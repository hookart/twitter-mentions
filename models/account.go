package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	PublicKey     string
	TwitterHandle string
	ENS           string
	Verified      bool
}
