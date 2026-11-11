package db

import (
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Each user can write one Markdown article per game, independently of UserGame.
type UserGameArticle struct {
	UserID      int       `gorm:"primaryKey;autoIncrement:false;index" json:"userId"`
	GameErogsID int       `gorm:"primaryKey;autoIncrement:false;index" json:"gameErogsId"`
	Title       string    `gorm:"not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	User      *User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	GameErogs *GameErogs `gorm:"foreignKey:GameErogsID;references:ID" json:"gameErogs,omitempty"`
}

func CreateUserGameArticle(db *gorm.DB, userID, gameErogsID int, title, content string) error {
	if _, err := GetGameErogsByID(db, gameErogsID); err != nil {
		return err
	}
	article := UserGameArticle{
		UserID:      userID,
		GameErogsID: gameErogsID,
		Title:       strings.TrimSpace(title),
		Content:     content,
	}
	if err := db.Create(&article).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return err
	}
	return nil
}

func GetAllUserGameArticles(db *gorm.DB) ([]UserGameArticle, error) {
	var articles []UserGameArticle
	err := db.Model(&UserGameArticle{}).
		Preload("User").
		Preload("GameErogs").
		Order("user_game_articles.updated_at DESC, user_game_articles.user_id, user_game_articles.game_erogs_id").
		Find(&articles).Error
	return articles, err
}

func GetUserGameArticlesByUserID(db *gorm.DB, userID int) ([]UserGameArticle, error) {
	if userID <= 0 {
		return nil, ErrParameterNotFound
	}
	var articles []UserGameArticle
	err := db.Model(&UserGameArticle{}).
		Preload("User").
		Preload("GameErogs").
		Where("user_game_articles.user_id = ?", userID).
		Order("user_game_articles.updated_at DESC, user_game_articles.user_id, user_game_articles.game_erogs_id").
		Find(&articles).Error
	return articles, err
}

func GetUserGameArticlesByGameErogsID(db *gorm.DB, gameErogsID int) ([]UserGameArticle, error) {
	if gameErogsID <= 0 {
		return nil, ErrParameterNotFound
	}
	var articles []UserGameArticle
	err := db.Model(&UserGameArticle{}).
		Preload("User").
		Preload("GameErogs").
		Where("user_game_articles.game_erogs_id = ?", gameErogsID).
		Order("user_game_articles.updated_at DESC, user_game_articles.user_id, user_game_articles.game_erogs_id").
		Find(&articles).Error
	return articles, err
}

func UpdateUserGameArticle(db *gorm.DB, userID, gameErogsID int, title, content string) error {
	res := db.Model(&UserGameArticle{}).
		Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).
		Select("title", "content", "updated_at").
		Updates(UserGameArticle{Title: strings.TrimSpace(title), Content: content, UpdatedAt: time.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteUserGameArticle(db *gorm.DB, userID, gameErogsID int) error {
	res := db.Where("user_id = ? AND game_erogs_id = ?", userID, gameErogsID).Delete(&UserGameArticle{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
