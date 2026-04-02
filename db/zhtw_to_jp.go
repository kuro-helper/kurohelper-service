package db

import "gorm.io/gorm"

// 撈出日文漢字以及繁體中文字對應資料
func GetAllZhtwToJps(db *gorm.DB) ([]ZhtwToJp, error) {
	var results []ZhtwToJp

	err := db.Find(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}
