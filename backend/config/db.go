package config

import (
	"easyinvesting/pkg/models"
	"errors"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbString := os.Getenv("EASYINVESTING_DB_STRING")
	isProd := os.Getenv("EASYINVESTING_PROD") == "true"
	if !isProd && dbString == "" {
		dbString = "db/easyinvesting.db"
	}
	if dbString == "" {
		panic(errors.New("EASYINVESTING_DB_STRING is not set"))
	}

	DB, err = gorm.Open(sqlite.Open(dbString), &gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	if err = DB.AutoMigrate(
		models.AllModels...,
	); err != nil {
		panic(err)
	}
}
