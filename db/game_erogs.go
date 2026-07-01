package db

import (
	"time"

	"gorm.io/gorm"
)

type GameErogs struct {
	ID           int       `gorm:"primaryKey;autoIncrement:false" json:"id"`
	BrandErogsID int       `json:"brandErogsId"`
	Name         string    `gorm:"not null" json:"name"` // 遊戲名稱(批評空間)
	Image        string    `gorm:"not null;default:''" json:"image"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	BrandErogs *BrandErogs `gorm:"foreignKey:BrandErogsID;references:ID" json:"brandErogs,omitempty"` // 單向 preload
}

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

func UpdateGameErogs(db *gorm.DB, id int, game *GameErogs) error {
	game.UpdatedAt = time.Now()
	return db.Model(&GameErogs{}).Where("id = ?", id).
		Select("Name", "BrandErogsID", "Image", "UpdatedAt").
		Updates(game).Error
}

func GetAllGameErogs(db *gorm.DB) ([]GameErogs, error) {
	var games []GameErogs
	err := db.Preload("BrandErogs").Find(&games).Error
	return games, err
}

func GetGameErogsByID(db *gorm.DB, id int) (GameErogs, error) {
	var game GameErogs
	err := db.First(&game, id).Error
	return game, err
}
