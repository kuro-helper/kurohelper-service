package db

type BrandCount struct {
	BrandID   int
	BrandName string
	Count     int
}

func GetUserHasPlayedBrandCount(userID string) ([]BrandCount, error) {
	var result []BrandCount

	err := Dbs.
		Table("user_games AS ug").
		Select(`
			brand.id AS brand_id,
			brand.name AS brand_name,
			COUNT(ug.game_erogs_id) AS count
		`).
		Joins("JOIN users u ON u.id = ug.user_id").
		Joins("JOIN game_erogs game ON game.id = ug.game_erogs_id").
		Joins("JOIN brand_erogs brand ON brand.id = game.brand_erogs_id").
		Where("u.discord_id = ?", userID).
		Where("ug.status = ?", 1).
		Group("brand.id, brand.name").
		Order("count DESC").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}
