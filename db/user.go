package db

import "gorm.io/gorm"

func EnsureUser(db *gorm.DB, userID, userName string) (*User, error) {
	var user User
	if err := db.Where("id = ?", userID).FirstOrCreate(&user, User{ID: userID, Name: userName}).Error; err != nil {
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

// 取得所有使用者資料
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var user []User

	err := db.Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
