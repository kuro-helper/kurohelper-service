package db

import (
	"time"

	"gorm.io/gorm"
)

// Game 自家遊戲主檔（目前僅作 erogs 無圖時的備選圖；未來再擴成主表）
type Game struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	ErogsID     *int   `gorm:"uniqueIndex" json:"erogsId"` // nil = 無對應批評空間
	ImageURL    string `gorm:"not null;default:''" json:"imageUrl"`
	UpdatedUser int    `gorm:"not null;default:0" json:"updatedUser"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// 以下為未來主表欄位，暫不啟用
	// Name        string
	// PlayHours   float64
	// Brand       int // preload
	// Links       int // preload
	// Description string
	// ScoreAvg    float64
	// ScoreCount  int
	// ReleaseDate time.Time
	// CreatedUser int
	// BangumiID   int
	// VNDBID      int `gorm:"column:vndb_id"`
	// Alias       int // preload
	// Tag         int // preload
}

func CreateGame(db *gorm.DB, erogsID *int, imageURL string, updatedUser int) (Game, error) {
	item := Game{
		ErogsID:     erogsID,
		ImageURL:    imageURL,
		UpdatedUser: updatedUser,
	}
	err := db.Create(&item).Error
	return item, err
}

func GetGameByID(db *gorm.DB, id int) (Game, error) {
	var item Game
	err := db.First(&item, id).Error
	return item, err
}

func GetGameByErogsID(db *gorm.DB, erogsID int) (Game, error) {
	var item Game
	err := db.Where("erogs_id = ?", erogsID).First(&item).Error
	return item, err
}

func UpdateGameImageURL(db *gorm.DB, id int, imageURL string, updatedUser int) error {
	res := db.Model(&Game{}).
		Where("id = ?", id).
		Updates(Game{
			ImageURL:    imageURL,
			UpdatedUser: updatedUser,
			UpdatedAt:   time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func UpdateGameImageURLByErogsID(db *gorm.DB, erogsID int, imageURL string, updatedUser int) error {
	res := db.Model(&Game{}).
		Where("erogs_id = ?", erogsID).
		Updates(Game{
			ImageURL:    imageURL,
			UpdatedUser: updatedUser,
			UpdatedAt:   time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func DeleteGame(db *gorm.DB, id int) error {
	res := db.Delete(&Game{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func DeleteGameByErogsID(db *gorm.DB, erogsID int) error {
	res := db.Where("erogs_id = ?", erogsID).Delete(&Game{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
