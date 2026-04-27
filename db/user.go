package db

import (
	"time"

	"gorm.io/gorm"
)

type (
	//  User主要資料
	User struct {
		ID          int       `gorm:"primaryKey" json:"id"`
		Name        string    `gorm:"not null;default:''" json:"name"` // 顯示用的名稱
		DiscordID   string    `gorm:"not null;default:''" json:"discordId"`
		Avatar      string    `gorm:"not null;default:''" json:"avatar"`
		Description string    `gorm:"not null;default:''" json:"description"`
		Role        int       `gorm:"not null;default:0" json:"role"`
		CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

		Auth      *UserAuth  `gorm:"foreignKey:UserID" json:"-"`
		UserGames []UserGame `gorm:"foreignKey:UserID" json:"userGames"`
	}

	// User帳號資料
	UserAuth struct {
		UserID   int    `gorm:"primaryKey;autoIncrement:false" json:"userId"`
		Username string `gorm:"uniqueIndex;size:30;not null" json:"username"`
		Password string `gorm:"not null" json:"-"`

		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
		UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	}
)

func EnsureDiscordUser(db *gorm.DB, discordID, userName string) (*User, error) {
	var user User
	if err := db.Where("discord_id = ?", discordID).FirstOrCreate(&user, User{DiscordID: discordID, Name: userName}).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 依據 userID 取得單一使用者資料
func GetUser(db *gorm.DB, userID string) (User, error) {
	var user User

	err := db.First(&user, "id = ?", userID).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// 依據 discordID 取得單一使用者資料
func GetUserByDiscordID(db *gorm.DB, discordID string) (User, error) {
	var user User

	err := db.Where("discord_id = ?", discordID).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// 取得所有使用者資料
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var user []User

	err := db.Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
