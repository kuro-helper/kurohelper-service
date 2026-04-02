package db

import (
	"time"

	"gorm.io/gorm"
)

// 取出所有的web api token
func GetWebAPIToken(db *gorm.DB) ([]WebAPIToken, error) {
	var tokens []WebAPIToken
	if err := db.Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// expiresDuration為Token的有效時間，無期限expiresDuration傳0
func CreateWebAPIToken(db *gorm.DB, id string, expiresDuration time.Duration) error {
	var expiresAt *time.Time

	if expiresDuration > 0 {
		t := time.Now().Add(expiresDuration)
		expiresAt = &t
	}

	token := &WebAPIToken{
		ID:        id,
		ExpiresAt: expiresAt,
	}

	if err := db.Create(token).Error; err != nil {
		return err
	}

	return nil
}
