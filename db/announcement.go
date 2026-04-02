package db

import (
	"time"

	"gorm.io/gorm"
)

func CreateAnnouncement(db *gorm.DB, category string, content string, thumbnail *string, image *string) error {
	newUpdate := Announcement{
		Category:  category,
		Content:   content,
		Thumbnail: thumbnail,
		Image:     image,
	}
	return db.Create(&newUpdate).Error
}

func GetAllAnnouncements(db *gorm.DB) ([]Announcement, error) {
	var list []Announcement

	err := db.Order("created_at DESC").Find(&list).Error

	return list, err
}

func GetAnnouncementByID(db *gorm.DB, id int) (Announcement, error) {
	var item Announcement
	err := db.First(&item, id).Error
	return item, err
}

func UpdateAnnouncement(db *gorm.DB, id int, updateData Announcement) error {
	updateData.UpdatedAt = time.Now()
	return db.Model(&Announcement{}).
		Where("id = ?", id).
		Select("category", "content", "thumbnail", "image", "updated_at").
		Updates(&updateData).
		Error
}

func DeleteAnnouncement(db *gorm.DB, id int) error {
	return db.Delete(&Announcement{}, id).Error
}
