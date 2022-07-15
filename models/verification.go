package models

import "gorm.io/gorm"

type Verification struct {
	gorm.Model
	AccountID          uint
	Account            *Account
	VerificationString string
}
