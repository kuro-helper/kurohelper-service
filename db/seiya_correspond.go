package db

import (
	"time"

	"gorm.io/gorm"
)

// 誠也對應表，專門針對極端狀況去對應
type SeiyaCorrespond struct {
	GameName  string    `gorm:"primaryKey"`
	SeiyaURL  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// 撈出誠也對應資料
func GetAllSeiyaCorresponds(db *gorm.DB) ([]SeiyaCorrespond, error) {
	var results []SeiyaCorrespond

	err := db.Find(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}
