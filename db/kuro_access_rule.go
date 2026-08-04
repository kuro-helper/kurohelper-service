package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	KuroAccessScopeChannel = "channel"
	KuroAccessScopeGuild   = "guild"
)

// KuroAccessRule stores runtime overrides for Discord channel and guild access.
// The environment allow-list remains the fallback when no channel override exists.
type KuroAccessRule struct {
	ScopeType string    `gorm:"primaryKey;column:scope_type;size:16" json:"scopeType"`
	ScopeID   string    `gorm:"primaryKey;column:scope_id;size:32" json:"scopeId"`
	Enabled   bool      `gorm:"column:enabled;not null" json:"enabled"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func ListKuroAccessRules(database *gorm.DB) ([]KuroAccessRule, error) {
	var rules []KuroAccessRule
	err := database.Order("scope_type ASC, scope_id ASC").Find(&rules).Error
	return rules, err
}

func UpsertKuroAccessRule(database *gorm.DB, scopeType, scopeID string, enabled bool) error {
	rule := KuroAccessRule{
		ScopeType: strings.ToLower(strings.TrimSpace(scopeType)),
		ScopeID:   strings.TrimSpace(scopeID),
		Enabled:   enabled,
		UpdatedAt: time.Now(),
	}
	if (rule.ScopeType != KuroAccessScopeChannel && rule.ScopeType != KuroAccessScopeGuild) || rule.ScopeID == "" {
		return ErrParameterNotFound
	}
	return database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "scope_type"}, {Name: "scope_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"enabled", "updated_at"}),
	}).Create(&rule).Error
}

func DeleteKuroAccessRule(database *gorm.DB, scopeType, scopeID string) error {
	scopeType = strings.ToLower(strings.TrimSpace(scopeType))
	scopeID = strings.TrimSpace(scopeID)
	if (scopeType != KuroAccessScopeChannel && scopeType != KuroAccessScopeGuild) || scopeID == "" {
		return ErrParameterNotFound
	}
	return database.Where("scope_type = ? AND scope_id = ?", scopeType, scopeID).Delete(&KuroAccessRule{}).Error
}
