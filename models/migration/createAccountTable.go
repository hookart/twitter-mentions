package migration

import "gorm.io/gorm"

func CreateAccountTable(db *gorm.DB) {
	type account struct {
		gorm.Model
		PublicKey     string `gorm:index:"idx_name,unique"`
		TwitterHandle string
		Verified      bool
		ENS           string
	}
	db.AutoMigrate(&account{})
}
