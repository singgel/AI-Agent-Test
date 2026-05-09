package models

import (
	"go-user-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	var err error
	
	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := autoMigrate(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(&User{})
}
