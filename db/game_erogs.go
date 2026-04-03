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

func UpdateGameErogsImage(db *gorm.DB, id int, game *GameErogs) error {
	game.UpdatedAt = time.Now()
	return db.Model(&GameErogs{}).Where("id = ?", id).
		Select("Image", "UpdatedAt").
		Updates(game).Error
}

func GetAllGameErogs(db *gorm.DB) ([]GameErogs, error) {
	var games []GameErogs
	err := db.Preload("BrandErogs").Find(&games).Error
	return games, err
}
