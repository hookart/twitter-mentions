package migration

import "gorm.io/gorm"

func CreateVerificationsTable(db *gorm.DB) {
	type verification struct {
		gorm.Model
		AccountID uint
		// Account            *Account
		VerificationString string
	}

	db.AutoMigrate(&verification{})
}
