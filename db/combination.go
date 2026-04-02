package db

type BrandCount struct {
	BrandID   int
	BrandName string
	Count     int
}

func GetUserHasPlayedBrandCount(userID string) ([]BrandCount, error) {
	var result []BrandCount

	err := Dbs.
		Table("user_has_playeds AS uhp").
		Select(`
			brand.id AS brand_id,
			brand.name AS brand_name,
			COUNT(uhp.game_erogs_id) AS count
		`).
		Joins("JOIN game_erogs game ON game.id = uhp.game_erogs_id").
		Joins("JOIN brand_erogs brand ON brand.id = game.brand_erogs_id").
		Where("uhp.user_id = ?", userID).
		Group("brand.id, brand.name").
		Order("count DESC").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}
