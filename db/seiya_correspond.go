package db

import "gorm.io/gorm"

// 撈出誠也對應資料
func GetAllSeiyaCorresponds(db *gorm.DB) ([]SeiyaCorrespond, error) {
	var results []SeiyaCorrespond

	err := db.Find(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}
