package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Title struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null;default:''" json:"description"`
	Symbol      string    `gorm:"not null;default:''" json:"symbol"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func CreateTitle(db *gorm.DB, name, description, symbol string) (Title, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Title{}, ErrParameterNotFound
	}

	item := Title{
		Name:        name,
		Description: strings.TrimSpace(description),
		Symbol:      strings.TrimSpace(symbol),
	}
	err := db.Create(&item).Error
	return item, err
}

func GetTitle(db *gorm.DB, id int) (Title, error) {
	var item Title
	err := db.First(&item, id).Error
	return item, err
}

func GetAllTitles(db *gorm.DB) ([]Title, error) {
	var list []Title
	err := db.Order("created_at DESC").Find(&list).Error
	return list, err
}

func DeleteTitle(db *gorm.DB, id int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Delete links first in case ON DELETE CASCADE is missing.
		if err := tx.Where("title_id = ?", id).Delete(&UserTitle{}).Error; err != nil {
			return err
		}
		res := tx.Delete(&Title{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
