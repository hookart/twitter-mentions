package migration

import "github.com/hookart/twitter-mentions/models"

func Migrate() {
	db := models.GetDBConnection()
	CreateAccountTable(db)
	CreateVerificationsTable(db)
}
