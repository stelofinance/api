package database

import (
	"github.com/stelofinance/api/models"
	"github.com/stelofinance/api/tools"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func ConnectDb() error {
	dsn, err := tools.GetEnvVariable("DB_CONNECTION_STRING")
	if err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	if prodEnv, _ := tools.GetEnvVariable("PRODUCTION_ENV"); prodEnv != "true" {
		db.Logger = logger.Default.LogMode(logger.Info)
	}

	Db = db
	return nil
}

func AutoMigrate() error {
	err := Db.AutoMigrate(&models.User{}, &models.Wallet{})
	return err
}
