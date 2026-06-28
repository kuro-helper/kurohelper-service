package db

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	UserRoleUser      = 0  // 一般使用者
	UserRoleDeveloper = 5  // 開發者
	UserRoleOwner     = 10 // 站主
)

func ValidUserRole(role int) bool {
	return role >= UserRoleUser && role <= UserRoleOwner
}

type User struct {
	ID              int       `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"not null;default:''" json:"name"` // 顯示用的名稱(暱稱)
	DiscordID       *string   `gorm:"uniqueIndex" json:"discordId"`
	Avatar          string    `gorm:"not null;default:''" json:"avatar"`
	Description     string    `gorm:"not null;default:''" json:"description"`
	PrivateGameData bool      `gorm:"not null;default:false" json:"privateGameData"` // 個人建檔隱私資料
	Role            int       `gorm:"not null;default:0" json:"role"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Auth      *UserAuth  `gorm:"foreignKey:UserID" json:"-"`
	UserGames []UserGame `gorm:"foreignKey:UserID" json:"userGames"`
}

func EnsureDiscordUser(db *gorm.DB, discordID, userName string) (*User, error) {
	if discordID == "" {
		return nil, ErrParameterNotFound
	}

	var user User
	if err := db.Where("discord_id = ?", discordID).FirstOrCreate(&user, User{DiscordID: &discordID, Name: userName}).Error; err != nil {
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

// 依據 discordID 取得單一使用者資料(含Auth)
func GetUserWithAuthByDiscordID(db *gorm.DB, discordID string) (User, error) {
	var user User

	err := db.
		Preload("Auth").
		Where("discord_id = ?", discordID).
		First(&user).Error
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

// 依據 discordID 更新 PrivateGameData 欄位
func UpdateUserPrivateGameDataByDiscordID(db *gorm.DB, discordID string, privateGameData bool) error {
	return db.Model(&User{}).
		Where("discord_id = ?", discordID).
		Update("private_game_data", privateGameData).Error
}

// 更新使用者個人資料（名稱、說明、大頭照 URL、建檔隱私）
func UpdateUser(db *gorm.DB, userID int, name, description, avatar string, privateGameData bool) (User, error) {
	var user User
	err := db.Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"name":              name,
			"description":       description,
			"avatar":            avatar,
			"private_game_data": privateGameData,
		}).Error
	if err != nil {
		return user, err
	}
	return GetUser(db, strconv.Itoa(userID))
}
