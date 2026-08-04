package kuro

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"kurohelperservice/db"
)

type AIStats = db.KuroAIStats
type AIProviderStats = db.KuroAIProviderStats
type AccessRule = db.KuroAccessRule

const (
	AccessScopeChannel = db.KuroAccessScopeChannel
	AccessScopeGuild   = db.KuroAccessScopeGuild
)

func ListAccessRules() ([]AccessRule, error) {
	return db.ListKuroAccessRules(db.Dbs)
}

func SetAccessRule(scopeType, scopeID string, enabled bool) error {
	return db.UpsertKuroAccessRule(db.Dbs, scopeType, scopeID, enabled)
}

func DeleteAccessRule(scopeType, scopeID string) error {
	return db.DeleteKuroAccessRule(db.Dbs, scopeType, scopeID)
}

func GetContextBoundary(channelID string) (string, error) {
	state, err := db.GetKuroChannelState(db.Dbs, channelID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return state.ContextBoundary, nil
}

func SetContextBoundary(channelID, messageID string) error {
	return db.SetKuroContextBoundary(db.Dbs, channelID, messageID)
}

func GetContextMessageKinds(channelID string, messageIDs []string) (map[string]string, error) {
	return db.GetKuroContextMessageKinds(db.Dbs, channelID, messageIDs)
}

func RecordCommandResponse(channelID, messageID string) error {
	return db.RecordKuroContextMessage(
		db.Dbs,
		channelID,
		messageID,
		db.KuroContextMessageCommandResponse,
	)
}

func DeleteContextMessagesCreatedBefore(cutoff time.Time) (int64, error) {
	return db.DeleteKuroContextMessagesCreatedBefore(db.Dbs, cutoff)
}

func GetAIStats(since time.Time) (AIStats, []AIProviderStats, error) {
	stats, err := db.GetKuroAIStats(db.Dbs, since)
	if err != nil {
		return AIStats{}, nil, err
	}
	providers, err := db.GetKuroAIProviderStats(db.Dbs, since)
	if err != nil {
		return AIStats{}, nil, err
	}
	return stats, providers, nil
}
