package db

import (
	"time"

	"gorm.io/gorm"
)

// 確保指定的BrandErogs存在，不存在就直接建立
func EnsureBrandErogs(db *gorm.DB, brandID int, brandName string, disband bool, gameCount int) (*BrandErogs, error) {
	var brand BrandErogs
	if err := db.Where("id = ?", brandID).FirstOrCreate(&brand, BrandErogs{ID: brandID, Name: brandName, Disband: disband, GameCount: gameCount}).Error; err != nil {
		return nil, err
	}
	return &brand, nil
}

func UpdateBrandErogs(db *gorm.DB, id int, brand *BrandErogs) error {
	brand.UpdatedAt = time.Now()
	return db.Model(&BrandErogs{}).Where("id = ?", id).
		Select("Name", "Disband", "GameCount", "UpdatedAt").
		Updates(brand).Error
}

func GetAllBrandErogs(db *gorm.DB) ([]BrandErogs, error) {
	var brands []BrandErogs
	err := db.Find(&brands).Error
	return brands, err
}
