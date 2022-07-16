package database

import (
	"github.com/stelofinance/api/tools"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func ConnectDb() error {
	dsn, err := tools.GetEnvVariable("DB_CONNECTION_STRING")
	if err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if isDev, _ := tools.GetEnvVariable("DEV_ENV"); isDev == "true" {
		db.Logger = logger.Default.LogMode(logger.Info)
	}

	// TODO: add migrations and models

	Db = db
	return nil
}
