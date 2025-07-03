package config

import (
	"errors"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
	if db != nil {
		return db // DB is already initialized
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

	if isProd {
		db, err = gorm.Open(postgres.Open(dbString), &gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			TranslateError:                           true,
		})
	} else {
		db, err = gorm.Open(sqlite.Open(dbString), &gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			TranslateError:                           true,
		})
	}

	if err != nil {
		log.Printf("failed to connect to database: %s", dbString)
		panic(err)
	}

	return db
}

