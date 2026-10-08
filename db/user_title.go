package db

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserTitle struct {
	UserID    int       `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	TitleID   int       `gorm:"primaryKey;autoIncrement:false" json:"titleId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

	User  *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Title *Title `gorm:"foreignKey:TitleID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func AddUserTitle(db *gorm.DB, userID, titleID int) error {
	if userID <= 0 || titleID <= 0 {
		return ErrParameterNotFound
	}

	row := UserTitle{UserID: userID, TitleID: titleID}
	if err := db.Create(&row).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return err
	}
	return nil
}

func GetTitlesByUserID(db *gorm.DB, userID int) ([]Title, error) {
	var list []Title
	err := db.
		Model(&Title{}).
		Joins("JOIN user_titles ON user_titles.title_id = titles.id").
		Where("user_titles.user_id = ?", userID).
		Order("titles.created_at DESC").
		Find(&list).Error
	return list, err
}

func RemoveUserTitle(db *gorm.DB, userID, titleID int) error {
	res := db.Where("user_id = ? AND title_id = ?", userID, titleID).Delete(&UserTitle{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
