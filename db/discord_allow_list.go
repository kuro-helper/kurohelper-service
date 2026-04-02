package db

import (
	"gorm.io/gorm"
)

// 查詢白名單
//
// 參數為guild(群組)跟dm(私訊)
func GetDiscordAllowListByKind(db *gorm.DB, kind string) ([]DiscordAllowList, error) {
	var results []DiscordAllowList

	if kind != "guild" && kind != "dm" {
		return nil, ErrParameterNotFound
	}

	err := db.Where("kind = ?", kind).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
