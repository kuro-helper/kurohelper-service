package db

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost     string
	DBOwner    string
	DBPassword string
	DBName     string
	DBPort     string
	SSLMode    string
}

// 全域連線池
var Dbs *gorm.DB

// 初始化資料庫連線
func InitDsn(config Config) error {
	if config.DBHost == "" {
		config.DBHost = "localhost"
	}
	config.SSLMode = strings.ToLower(strings.TrimSpace(config.SSLMode))
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost,
		config.DBOwner,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		config.SSLMode,
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
