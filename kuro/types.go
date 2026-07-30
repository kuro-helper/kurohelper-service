package kuro

import "time"

const ProtocolVersion = 1

type MentionedUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type RecentMessage struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Content     string    `json:"content"`
	Assistant   bool      `json:"assistant"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GenerateRequest struct {
	RequestID           string          `json:"requestId"`
	ChannelID           string          `json:"channelId"`
	UserID              string          `json:"userId"`
	DisplayName         string          `json:"displayName"`
	Text                string          `json:"text"`
	RecentContext       string          `json:"recentChannelContext,omitempty"`
	RetrievalText       string          `json:"retrievalText,omitempty"`
	MentionedUsers      []MentionedUser `json:"mentionedUsers,omitempty"`
	ContextParticipants []MentionedUser `json:"contextParticipants,omitempty"`
}

type GenerateResponse struct {
	Text string `json:"text"`
}

type Memory struct {
	ID           string              `json:"id"`
	SubjectName  string              `json:"subject_name"`
	Key          string              `json:"key"`
	Value        string              `json:"value"`
	Category     string              `json:"category"`
	Status       string              `json:"status"`
	UpdatedAt    string              `json:"updated_at"`
	PurgeAfter   string              `json:"purge_after,omitempty"`
	Scope        string              `json:"scope,omitempty"`
	ScopeID      string              `json:"scope_id,omitempty"`
	Participants []MemoryParticipant `json:"participants,omitempty"`
}

type MemoryParticipant struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role,omitempty"`
}

type MemoryResponse struct {
	Status             string   `json:"status"`
	Count              int      `json:"count,omitempty"`
	TrashRetentionDays int      `json:"trash_retention_days,omitempty"`
	Memories           []Memory `json:"memories,omitempty"`
	Memory             *Memory  `json:"memory,omitempty"`
	ConflictID         string   `json:"conflict_id,omitempty"`
}

type MemoryRequest struct {
	Action   string `json:"action"`
	Status   string `json:"status,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	MemoryID string `json:"memoryId,omitempty"`
}

type HealthResponse struct {
	Status             string `json:"status"`
	SillyTavernReady   bool   `json:"sillyTavernReady"`
	MemoryEnabled      bool   `json:"memoryEnabled"`
	TrashRetentionDays int    `json:"trashRetentionDays"`
}
