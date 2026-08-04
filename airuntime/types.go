package airuntime

import "time"

const ProtocolVersion = 1

type MentionedUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type RecentMessage struct {
	ID          string            `json:"id"`
	UserID      string            `json:"userId"`
	DisplayName string            `json:"displayName"`
	Content     string            `json:"content"`
	Assistant   bool              `json:"assistant"`
	CreatedAt   time.Time         `json:"createdAt"`
	Images      []ImageAttachment `json:"images,omitempty"`
}

type ImageAttachment struct {
	ID                string `json:"id,omitempty"`
	URL               string `json:"url"`
	Filename          string `json:"filename,omitempty"`
	ContentType       string `json:"contentType,omitempty"`
	Size              int    `json:"size,omitempty"`
	MessageID         string `json:"messageId,omitempty"`
	AuthorName        string `json:"authorName,omitempty"`
	SourceKind        string `json:"sourceKind,omitempty"`
	SourceMessageText string `json:"sourceMessageText,omitempty"`
	ContextOnly       bool   `json:"contextOnly,omitempty"`
}

type GenerateRequest struct {
	RequestID           string            `json:"requestId"`
	ChannelID           string            `json:"channelId"`
	UserID              string            `json:"userId"`
	DisplayName         string            `json:"displayName"`
	Text                string            `json:"text"`
	RecentContext       string            `json:"recentChannelContext,omitempty"`
	RecentMessages      []RecentMessage   `json:"recentMessages,omitempty"`
	RetrievalText       string            `json:"retrievalText,omitempty"`
	MentionedUsers      []MentionedUser   `json:"mentionedUsers,omitempty"`
	ContextParticipants []MentionedUser   `json:"contextParticipants,omitempty"`
	Images              []ImageAttachment `json:"images,omitempty"`
}

type GenerateResponse struct {
	Text    string             `json:"text"`
	Metrics *GenerationMetrics `json:"metrics,omitempty"`
}

type MetricEvent struct {
	RequestID string             `json:"requestId"`
	ChannelID string             `json:"channelId"`
	Operation string             `json:"operation"`
	Metrics   *GenerationMetrics `json:"metrics,omitempty"`
}

type GenerationMetrics struct {
	Status                  string   `json:"status,omitempty"`
	Model                   string   `json:"model,omitempty"`
	Provider                string   `json:"provider,omitempty"`
	Providers               []string `json:"providers,omitempty"`
	ProviderModel           string   `json:"providerModel,omitempty"`
	RoutingStrategy         string   `json:"routingStrategy,omitempty"`
	RoutingRegion           string   `json:"routingRegion,omitempty"`
	RoutingAttempt          int      `json:"routingAttempt,omitempty"`
	GenerationIDs           []string `json:"generationIds,omitempty"`
	ProviderStatusCode      int      `json:"providerStatusCode,omitempty"`
	UsageAvailable          bool     `json:"usageAvailable"`
	PromptTokens            int64    `json:"promptTokens,omitempty"`
	CompletionTokens        int64    `json:"completionTokens,omitempty"`
	TotalTokens             int64    `json:"totalTokens,omitempty"`
	ReasoningTokens         int64    `json:"reasoningTokens,omitempty"`
	CachedTokens            int64    `json:"cachedTokens,omitempty"`
	CostUSD                 float64  `json:"costUsd,omitempty"`
	GenerationCount         int      `json:"generationCount,omitempty"`
	RetryCount              int      `json:"retryCount,omitempty"`
	MemoryRecallMs          int64    `json:"memoryRecallMs,omitempty"`
	BrowserQueueMs          int64    `json:"browserQueueMs,omitempty"`
	LocaleMs                int64    `json:"localeMs,omitempty"`
	PersonaMs               int64    `json:"personaMs,omitempty"`
	InjectUserMs            int64    `json:"injectUserMs,omitempty"`
	GenerateMs              int64    `json:"generateMs,omitempty"`
	ValidateMs              int64    `json:"validateMs,omitempty"`
	RetryGenerateMs         int64    `json:"retryGenerateMs,omitempty"`
	RetryValidateMs         int64    `json:"retryValidateMs,omitempty"`
	BrowserTotalMs          int64    `json:"browserTotalMs,omitempty"`
	EndToEndBeforeDiscordMs int64    `json:"endToEndBeforeDiscordMs,omitempty"`
	ProviderHeadersMs       int64    `json:"providerHeadersMs,omitempty"`
	ProviderFirstTokenMs    int64    `json:"providerFirstTokenMs,omitempty"`
	ProviderDurationMs      int64    `json:"providerDurationMs,omitempty"`
}

type Memory struct {
	ID              string              `json:"id"`
	SubjectID       string              `json:"subject_id,omitempty"`
	SubjectName     string              `json:"subject_name"`
	Key             string              `json:"key"`
	Value           string              `json:"value"`
	Category        string              `json:"category"`
	Importance      float64             `json:"importance,omitempty"`
	Confidence      float64             `json:"confidence,omitempty"`
	Status          string              `json:"status"`
	CreatedAt       string              `json:"created_at,omitempty"`
	UpdatedAt       string              `json:"updated_at"`
	LastAccessedAt  string              `json:"last_accessed_at,omitempty"`
	AccessCount     int                 `json:"access_count,omitempty"`
	PurgeAfter      string              `json:"purge_after,omitempty"`
	Scope           string              `json:"scope,omitempty"`
	ScopeID         string              `json:"scope_id,omitempty"`
	SupersedesID    string              `json:"supersedes_id,omitempty"`
	SourceRequestID string              `json:"source_request_id,omitempty"`
	SourceChannelID string              `json:"source_channel_id,omitempty"`
	Participants    []MemoryParticipant `json:"participants,omitempty"`
}

type MemoryParticipant struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role,omitempty"`
}

type MemoryBackup struct {
	ID           string         `json:"id"`
	CreatedAt    string         `json:"created_at"`
	Reason       string         `json:"reason"`
	SizeBytes    int64          `json:"size_bytes"`
	MemoryCount  int            `json:"memory_count"`
	StatusCounts map[string]int `json:"status_counts,omitempty"`
}

type MemoryResponse struct {
	Status               string         `json:"status"`
	Count                int            `json:"count,omitempty"`
	TrashRetentionDays   int            `json:"trash_retention_days,omitempty"`
	Memories             []Memory       `json:"memories,omitempty"`
	Memory               *Memory        `json:"memory,omitempty"`
	ConflictID           string         `json:"conflict_id,omitempty"`
	BackupRetentionCount int            `json:"backup_retention_count,omitempty"`
	Backups              []MemoryBackup `json:"backups,omitempty"`
	Backup               *MemoryBackup  `json:"backup,omitempty"`
	SafetyBackup         *MemoryBackup  `json:"safety_backup,omitempty"`
	RestoredActiveCount  int            `json:"restored_active_count,omitempty"`
}

type MemoryRequest struct {
	Action   string `json:"action"`
	Status   string `json:"status,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	MemoryID string `json:"memoryId,omitempty"`
	BackupID string `json:"backupId,omitempty"`
}

type HealthResponse struct {
	Status             string            `json:"status"`
	SillyTavernReady   bool              `json:"sillyTavernReady"`
	MemoryEnabled      bool              `json:"memoryEnabled"`
	TrashRetentionDays int               `json:"trashRetentionDays"`
	VisionCache        *VisionCacheStats `json:"visionCache,omitempty"`
}

type VisionCacheStats struct {
	Enabled                 bool    `json:"enabled"`
	Persistent              bool    `json:"persistent"`
	Entries                 int     `json:"entries"`
	Images                  int     `json:"images"`
	OCREntries              int     `json:"ocrEntries"`
	ObservationEntries      int     `json:"observationEntries"`
	MaxImages               int     `json:"maxImages"`
	MaxObservationsPerImage int     `json:"maxObservationsPerImage"`
	Hits                    int64   `json:"hits"`
	Misses                  int64   `json:"misses"`
	HitRate                 float64 `json:"hitRate"`
	Writes                  int64   `json:"writes"`
	LoadedEntries           int64   `json:"loadedEntries"`
	ExpiredEntries          int64   `json:"expiredEntries"`
	EvictedImages           int64   `json:"evictedImages"`
	EvictedEntries          int64   `json:"evictedEntries"`
	EvictedObservations     int64   `json:"evictedObservations"`
}

type RawReply struct {
	CachedAt string `json:"cachedAt,omitempty"`
	RawText  string `json:"rawText"`
	Source   string `json:"source,omitempty"`
}

type RawRepliesResponse struct {
	Entries []RawReply `json:"entries"`
}
