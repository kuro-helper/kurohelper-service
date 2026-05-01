package db

import (
	"time"

	"gorm.io/gorm"
)

type RegisterCache struct {
	ID        string    `gorm:"primaryKey"`
	DiscordID string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func CreateRegisterCache(db *gorm.DB, id, discordID string, expiresDuration time.Duration) error {
	cache := RegisterCache{
		ID:        id,
		DiscordID: discordID,
		ExpiresAt: time.Now().Add(expiresDuration),
	}

	if err := db.Create(&cache).Error; err != nil {
		return err
	}

	return nil
}

// 取得尚未過期的註冊快取資料。
func GetRegisterCacheByID(db *gorm.DB, id string) (RegisterCache, error) {
	var cache RegisterCache

	err := db.
		Where("id = ?", id).
		Where("expires_at > ?", time.Now()).
		First(&cache).Error
	if err != nil {
		return cache, err
	}

	return cache, nil
}

func DeleteRegisterCacheByID(db *gorm.DB, id string) error {
	res := db.Where("id = ?", id).Delete(&RegisterCache{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
