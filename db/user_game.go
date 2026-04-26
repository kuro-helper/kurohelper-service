package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// User的遊戲資料
type UserGame struct {
	UserID      int `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	GameErogsID int `gorm:"primaryKey;autoIncrement:false" json:"gameErogsId"`
	// GameID string `gorm:"primaryKey" json:"gameId"`

	Status        int  `gorm:"not null;default:0" json:"status"`
	WishListMark  bool `gorm:"not null;default:false" json:"wishListMark"`
	BlackListMark bool `gorm:"not null;default:false" json:"blackListMark"`

	StartDate    *time.Time `json:"startDate,omitempty"`
	FinishedDate *time.Time `json:"finishedDate,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	GameErogs *GameErogs `gorm:"foreignKey:GameErogsID;references:ID" json:"gameErogs,omitempty"` // 單向 preload
	// Game *Game `gorm:"foreignKey:GameID;references:ID" json:"game"`
}

func GetUserGameFinishedByID(db *gorm.DB, userID int) ([]UserGame, error) {
	var hasPlayed []UserGame

	err := db.
		Model(&UserGame{}).
		Preload("GameErogs").
		// Preload("GameErogs.BrandErogs").
		Where("user_id = ?", userID).
		Where("status = ?", 1).
		Order("COALESCE(finished_date, created_at) DESC").
		Find(&hasPlayed).Error
	if err != nil {
		return nil, err
	}

	return hasPlayed, nil
}

// 有transaction，因為要先檢查DiscordID正確性，到時候再往上移
func GetUserGameByDiscordID(db *gorm.DB, discordID string) ([]UserGame, error) {
	var hasPlayed []UserGame
	err := db.Transaction(func(tx *gorm.DB) error {
		var users []User
		if err := tx.Where("discord_id = ?", discordID).Limit(2).Find(&users).Error; err != nil {
			return err
		}
		if len(users) == 0 {
			return gorm.ErrRecordNotFound
		}
		if len(users) > 1 {
			return fmt.Errorf("expected one user by discord_id %s, got %d", discordID, len(users))
		}

		return tx.
			Model(&UserGame{}).
			Preload("GameErogs").
			Where("user_id = ?", users[0].ID).
			Order("COALESCE(finished_date, created_at) DESC").
			Find(&hasPlayed).Error
	})
	if err != nil {
		return nil, err
	}

	return hasPlayed, nil
}

func GetUserGameByUserAndGameNameLike(db *gorm.DB, userID int, gameErogsName string) (UserGame, error) {
	var result UserGame

	err := db.
		Model(&UserGame{}).
		Joins("JOIN game_erogs ON game_erogs.id = user_games.game_erogs_id").
		Where("user_games.user_id = ?", userID).
		Where("game_erogs.name ILIKE ?", "%"+gameErogsName+"%").
		Preload("GameErogs").
		First(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

// EnsureUserGame ensures an empty user_games row exists.
// It only guarantees PK row existence and keeps other fields at defaults.
func EnsureUserGame(db *gorm.DB, userID int, gameErogsID int) error {
	userGame := UserGame{
		UserID:      userID,
		GameErogsID: gameErogsID,
	}

	if err := db.Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).FirstOrCreate(&userGame).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return err
	}

	return nil
}

func UpdateUserGameFinished(db *gorm.DB, userID int, gameErogsID int, completedAt *time.Time) error {
	res := db.Model(&UserGame{}).
		Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).
		Select("status", "finished_date", "updated_at").
		Updates(UserGame{
			Status:       1,
			FinishedDate: completedAt,
			UpdatedAt:    time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateUserGameWishListMark(db *gorm.DB, userID, gameErogsID int, wishListMark bool) error {
	res := db.Model(&UserGame{}).
		Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).
		Select("wish_list_mark", "updated_at").
		Updates(UserGame{
			WishListMark: wishListMark,
			UpdatedAt:    time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateUserGameBlackListMark(db *gorm.DB, userID, gameErogsID int, blackListMark bool) error {
	res := db.Model(&UserGame{}).
		Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).
		Select("black_list_mark", "updated_at").
		Updates(UserGame{
			BlackListMark: blackListMark,
			UpdatedAt:     time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 刪除 User 的遊戲資料
//
// 這邊是有transaction的版本，到時候再往上移
func DeleteUserGame(db *gorm.DB, discordID string, gameErogsID int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&User{}).Where("discord_id = ?", discordID).Count(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return fmt.Errorf("expected exactly one user by discord_id %s, got %d", discordID, count)
		}

		var user User
		if err := tx.Where("discord_id = ?", discordID).First(&user).Error; err != nil {
			return err
		}

		return tx.
			Where("user_id = ? AND game_erogs_id = ?", user.ID, gameErogsID).
			Delete(&UserGame{}).Error
	})
}
