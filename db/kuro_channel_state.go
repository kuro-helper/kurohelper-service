package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KuroChannelState struct {
	ChannelID       string    `gorm:"primaryKey;column:channel_id;size:32" json:"channelId"`
	ContextBoundary string    `gorm:"column:context_boundary;size:32;not null;default:''" json:"contextBoundary"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func GetKuroChannelState(db *gorm.DB, channelID string) (KuroChannelState, error) {
	var state KuroChannelState
	err := db.First(&state, "channel_id = ?", strings.TrimSpace(channelID)).Error
	return state, err
}

func SetKuroContextBoundary(db *gorm.DB, channelID, messageID string) error {
	state := KuroChannelState{
		ChannelID:       strings.TrimSpace(channelID),
		ContextBoundary: strings.TrimSpace(messageID),
		UpdatedAt:       time.Now(),
	}
	if state.ChannelID == "" || state.ContextBoundary == "" {
		return ErrParameterNotFound
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "channel_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"context_boundary", "updated_at"}),
	}).Create(&state).Error
}
