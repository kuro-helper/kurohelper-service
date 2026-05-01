package db

import (
	"time"

	"gorm.io/gorm"
)

// User帳號資料
type UserAuth struct {
	UserID   int    `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	Username string `gorm:"uniqueIndex;size:30;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// 依據 userName 取得單一使用者帳號資料
func GetUserAuthByUsername(db *gorm.DB, userName string) (UserAuth, error) {
	var auth UserAuth

	err := db.Where("username = ?", userName).First(&auth).Error
	if err != nil {
		return auth, err
	}

	return auth, nil
}

// 依據 userID 取得單一使用者帳號資料
func GetUserAuthByUserID(db *gorm.DB, userID int) (UserAuth, error) {
	var auth UserAuth

	err := db.Where("user_id = ?", userID).First(&auth).Error
	if err != nil {
		return auth, err
	}

	return auth, nil
}

// 建立使用者帳號資料
func CreateUserAuth(db *gorm.DB, userID int, userName, password string) error {
	auth := UserAuth{
		UserID:   userID,
		Username: userName,
		Password: password,
	}

	if err := db.Create(&auth).Error; err != nil {
		return err
	}

	return nil
}
