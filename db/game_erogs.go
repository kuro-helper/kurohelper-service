package db

import (
	"time"

	"gorm.io/gorm"
)

// 確保指定的GameErogs存在，不存在就直接建立
func EnsureGameErogs(db *gorm.DB, gameID int, gameName string, gameImage string, brandID int) (*GameErogs, error) {
	var game GameErogs
	if err := db.Where(GameErogs{ID: gameID}).
		Attrs(GameErogs{Name: gameName, BrandErogsID: brandID, Image: gameImage}).
		FirstOrCreate(&game).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func UpdateGameErogsImageByID(db *gorm.DB, id int, image string) error {
	return db.Model(&GameErogs{}).Where("id = ?", id).Updates(GameErogs{
		Image:     image,
		UpdatedAt: time.Now(),
	}).Error
}

func GetAllGameErogs(db *gorm.DB) ([]GameErogs, error) {
	var games []GameErogs
	err := db.Preload("BrandErogs").Find(&games).Error
	return games, err
}
