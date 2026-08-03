package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const KuroContextMessageCommandResponse = "command_response"

// KuroContextMessage stores Discord messages that have an operational role
// and must not be interpreted as conversation history by the model.
type KuroContextMessage struct {
	MessageID string    `gorm:"primaryKey;column:message_id;size:32" json:"messageId"`
	ChannelID string    `gorm:"index;column:channel_id;size:32;not null" json:"channelId"`
	Kind      string    `gorm:"index;column:kind;size:32;not null" json:"kind"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func RecordKuroContextMessage(database *gorm.DB, channelID, messageID, kind string) error {
	message := KuroContextMessage{
		MessageID: strings.TrimSpace(messageID),
		ChannelID: strings.TrimSpace(channelID),
		Kind:      strings.TrimSpace(kind),
	}
	if database == nil || message.MessageID == "" || message.ChannelID == "" || message.Kind == "" {
		return ErrParameterNotFound
	}
	return database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"channel_id", "kind"}),
	}).Create(&message).Error
}

func GetKuroContextMessageKinds(database *gorm.DB, channelID string, messageIDs []string) (map[string]string, error) {
	result := make(map[string]string)
	if database == nil || strings.TrimSpace(channelID) == "" || len(messageIDs) == 0 {
		return result, nil
	}
	var messages []KuroContextMessage
	err := database.Select("message_id", "kind").
		Where("channel_id = ? AND message_id IN ?", strings.TrimSpace(channelID), messageIDs).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		result[message.MessageID] = message.Kind
	}
	return result, nil
}

// DeleteKuroContextMessagesCreatedBefore removes expired operational message
// markers. It never deletes the corresponding Discord messages.
func DeleteKuroContextMessagesCreatedBefore(database *gorm.DB, cutoff time.Time) (int64, error) {
	if database == nil || cutoff.IsZero() {
		return 0, ErrParameterNotFound
	}
	result := database.
		Where("created_at < ?", cutoff).
		Delete(&KuroContextMessage{})
	return result.RowsAffected, result.Error
}
