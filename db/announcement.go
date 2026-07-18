package db

import (
	"time"

	"gorm.io/gorm"
)

// 公告資料表
type Announcement struct {
	ID        int       `gorm:"primaryKey"`
	Category  string    `gorm:"not null"`
	Title     string    `gorm:"not null"`  // 列表標題
	Content   string    `gorm:"not null"`  // markdown內文
	Icon      *string   // 網頁板小 icon（MDI 名稱，如 mdi-bullhorn）
	Thumbnail *string   // Discord側邊小圖（圖片 URL）
	Image     *string   // Discord底部大圖（圖片 URL）
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateAnnouncement(db *gorm.DB, category, title, content string, icon *string, thumbnail *string, image *string) (Announcement, error) {
	item := Announcement{
		Category:  category,
		Title:     title,
		Content:   content,
		Icon:      icon,
		Thumbnail: thumbnail,
		Image:     image,
	}
	err := db.Create(&item).Error
	return item, err
}

func GetAllAnnouncements(db *gorm.DB) ([]Announcement, error) {
	var list []Announcement
	err := db.Order("created_at DESC").Find(&list).Error
	return list, err
}

func GetAnnouncementsByCategory(db *gorm.DB, category string) ([]Announcement, error) {
	var list []Announcement
	err := db.Where("category = ?", category).Order("created_at DESC").Find(&list).Error
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
		Select("category", "title", "content", "icon", "thumbnail", "image", "updated_at").
		Updates(&updateData).
		Error
}

func DeleteAnnouncement(db *gorm.DB, id int) error {
	return db.Delete(&Announcement{}, id).Error
}
