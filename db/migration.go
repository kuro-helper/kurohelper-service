package db

import (
	"errors"

	"gorm.io/gorm"
)

// Migration
func Migration(db *gorm.DB) error {
	if db == nil {
		return errors.New("DB not initialized")
	}

	db.AutoMigrate(&ZhtwToJp{})
	db.AutoMigrate(&SeiyaCorrespond{})
	db.AutoMigrate(&WebAPIToken{})
	db.AutoMigrate(&DiscordAllowList{})
	db.AutoMigrate(
		&BrandErogs{},
		&GameErogs{},
	)
	db.AutoMigrate(
		&User{},
		&UserHasPlayed{},
		&UserInWish{},
	)

	// db.AutoMigrate(&Announcement{})

	return nil
}
