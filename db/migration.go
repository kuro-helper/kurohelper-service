package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Migration
func Migration(db *gorm.DB) error {
	if db == nil {
		return errors.New("DB not initialized")
	}

	if err := db.AutoMigrate(
		&ZhtwToJp{},
		&SeiyaCorrespond{},
		&WebAPIToken{},
		&RegisterCache{},
		&DiscordAllowList{},
		&BrandErogs{},
		&GameErogs{},
	); err != nil {
		return err
	}
	if err := dropErogsNameUniqueConstraints(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(
		&User{},
		&UserAuth{},
		&UserGame{},
		&UserCheckIn{},
		&AppConfig{},
		&Announcement{},
		&KuroChannelState{},
		&KuroAIMetric{},
		&KuroAIDailyMetric{},
	); err != nil {
		return err
	}

	// db.AutoMigrate(&History{})

	return nil
}

func dropErogsNameUniqueConstraints(db *gorm.DB) error {
	drops := []struct {
		table      string
		constraint string
	}{
		{"game_erogs", "uni_game_erogs_name"},
		{"brand_erogs", "uni_brand_erogs_name"},
		{"game_erogs", "game_erogs_name_key"},
		{"brand_erogs", "brand_erogs_name_key"},
	}
	for _, item := range drops {
		sql := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s", item.table, item.constraint)
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
