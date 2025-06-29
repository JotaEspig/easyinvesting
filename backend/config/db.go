package config

import (
	"errors"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	if DB != nil {
		return // DB is already initialized
	}

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
		TranslateError:                           true,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", dbString)
		panic(err)
	}
}
