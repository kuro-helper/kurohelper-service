package db

import (
	"time"

	"gorm.io/gorm"
)

type DiscordAllowList struct {
	ID         string    `gorm:"primaryKey"`
	Kind       string    `gorm:"not null"`
	Permission int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// 查詢白名單
//
// 參數為guild(群組)跟dm(私訊)
func GetDiscordAllowListByKind(db *gorm.DB, kind string) ([]DiscordAllowList, error) {
	var results []DiscordAllowList

	if kind != "guild" && kind != "dm" {
		return nil, ErrParameterNotFound
	}

	err := db.Where("kind = ?", kind).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
