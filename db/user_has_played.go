package db

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func GetUserHasPlayedByID(db *gorm.DB, userID string) ([]UserHasPlayed, error) {
	var hasPlayed []UserHasPlayed

	err := db.
		// Preload("User").
		Preload("GameErogs").
		Preload("GameErogs.BrandErogs").
		Where("user_id = ?", userID).
		Order("COALESCE(completed_at, created_at) DESC").
		Find(&hasPlayed).Error

	if err != nil {
		return nil, err
	}
	return hasPlayed, nil
}

func GetUserHasPlayedByUserAndGameID(db *gorm.DB, userID string, gameErogsID int) (UserHasPlayed, error) {
	var userHasPlayed UserHasPlayed

	err := db.First(&userHasPlayed, "user_id = ? AND game_erogs_id = ?", userID, gameErogsID).Error
	if err != nil {
		return userHasPlayed, err
	}

	return userHasPlayed, nil
}

func DeleteUserHasPlayed(db *gorm.DB, userID string, gameErogsID int) error {
	err := db.
		Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).
		Delete(&UserHasPlayed{}).Error

	return err
}

func CreateUserHasPlayed(db *gorm.DB, userID string, gameErogsID int, completedAt *time.Time) error {
	userHasPlayed := UserHasPlayed{
		UserID:      userID,
		GameErogsID: gameErogsID,
		CompletedAt: completedAt,
	}

	if err := db.Create(&userHasPlayed).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return err
	}

	return nil
}

func GetUserHasPlayedByUserAndGameNameLike(db *gorm.DB, userID string, gameErogsName string) (UserHasPlayed, error) {
	var result UserHasPlayed

	err := db.
		Model(&UserHasPlayed{}).
		Joins("JOIN game_erogs ON game_erogs.id = user_has_playeds.game_erogs_id").
		Where("user_has_playeds.user_id = ?", userID).
		Where("game_erogs.name ILIKE ?", "%"+gameErogsName+"%").
		Preload("GameErogs").
		First(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}
