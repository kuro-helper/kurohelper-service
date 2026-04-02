package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func GetUserInWishByID(db *gorm.DB, userID string) ([]UserInWish, error) {
	var inWish []UserInWish

	err := db.
		// Preload("User").
		Preload("GameErogs").
		Preload("GameErogs.BrandErogs").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&inWish).Error

	if err != nil {
		return nil, err
	}
	return inWish, nil
}

func GetUserInWishByUserAndGameID(db *gorm.DB, userID string, gameErogsID int) (UserInWish, error) {
	var userInWish UserInWish

	err := db.First(&userInWish, "user_id = ? AND game_erogs_id = ?", userID, gameErogsID).Error
	if err != nil {
		return userInWish, err
	}

	return userInWish, nil
}

func DeleteUserInWish(db *gorm.DB, userID string, gameErogsID int) error {
	err := db.
		Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).
		Delete(&UserInWish{}).Error

	return err
}

func CreateUserInWish(db *gorm.DB, userID string, gameErogsID int) error {
	userInWish := UserInWish{
		UserID:      userID,
		GameErogsID: gameErogsID,
	}

	if err := db.Create(&userInWish).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return err
	}

	return nil
}

func GetUserInWishByUserAndGameNameLike(db *gorm.DB, userID string, gameName string) (UserInWish, error) {
	var result UserInWish

	err := db.
		Model(&UserInWish{}).
		Joins("JOIN game_erogs ON game_erogs.id = user_in_wishes.game_erogs_id").
		Where("user_in_wishes.user_id = ?", userID).
		Where("game_erogs.name ILIKE ?", "%"+gameName+"%").
		Preload("GameErogs").
		First(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}
