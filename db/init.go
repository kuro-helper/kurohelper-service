package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBOwner    string
	DBPassword string
	DBName     string
	DBPort     string
}

// 全域連線池
var Dbs *gorm.DB

// 初始化資料庫連線
func InitDsn(config Config) error {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBOwner,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)

	// get connect db variable
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	Dbs = db

	return nil
}
