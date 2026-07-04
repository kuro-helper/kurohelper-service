package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// AppConfig 系統設定表
type AppConfig struct {
	ConfigKey   string    `gorm:"primaryKey;column:config_key;size:255" json:"configKey"`
	ConfigValue string    `gorm:"column:config_value;type:text;not null" json:"configValue"`
	Description string    `gorm:"size:255;not null;default:''" json:"description"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func CreateAppConfig(db *gorm.DB, configKey, configValue, description string) error {
	configKey = strings.TrimSpace(configKey)
	if configKey == "" {
		return ErrParameterNotFound
	}

	item := AppConfig{
		ConfigKey:   configKey,
		ConfigValue: configValue,
		Description: strings.TrimSpace(description),
		UpdatedAt:   time.Now(),
	}
	return db.Create(&item).Error
}

func GetAppConfigByKey(db *gorm.DB, configKey string) (AppConfig, error) {
	var item AppConfig
	err := db.First(&item, "config_key = ?", strings.TrimSpace(configKey)).Error
	return item, err
}

func GetAllAppConfigs(db *gorm.DB) ([]AppConfig, error) {
	var list []AppConfig
	err := db.Order("config_key ASC").Find(&list).Error
	return list, err
}

func UpdateAppConfig(db *gorm.DB, configKey, configValue, description string) error {
	configKey = strings.TrimSpace(configKey)
	if configKey == "" {
		return ErrParameterNotFound
	}

	res := db.Model(&AppConfig{}).
		Where("config_key = ?", configKey).
		Select("config_value", "description", "updated_at").
		Updates(AppConfig{
			ConfigValue: configValue,
			Description: strings.TrimSpace(description),
			UpdatedAt:   time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteAppConfigByKey(db *gorm.DB, configKey string) error {
	configKey = strings.TrimSpace(configKey)
	if configKey == "" {
		return ErrParameterNotFound
	}

	res := db.Where("config_key = ?", configKey).Delete(&AppConfig{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
