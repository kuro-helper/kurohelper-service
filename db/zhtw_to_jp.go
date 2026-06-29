package db

import (
	"time"

	"gorm.io/gorm"
)

type ZhtwToJp struct {
	ZhTw      string    `gorm:"primaryKey;size:1"` // 繁體中文漢字
	Jp        string    `gorm:"size:1;not null"`   // 日文漢字
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (ZhtwToJp) TableName() string {
	return "zhtw_to_jp"
}

// 撈出日文漢字以及繁體中文字對應資料
func GetAllZhtwToJps(db *gorm.DB) ([]ZhtwToJp, error) {
	var results []ZhtwToJp

	err := db.Find(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}
