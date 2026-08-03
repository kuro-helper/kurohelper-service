package kuro

import (
	"time"

	"kurohelperservice/airuntime"
	"kurohelperservice/db"
)

type GenerationMetricRecord struct {
	RequestID          string
	ChannelID          string
	CreatedAt          time.Time
	Status             string
	Metrics            *airuntime.GenerationMetrics
	BotQueueMs         int64
	DiscordHistoryMs   int64
	RuntimeRoundTripMs int64
	DiscordSendMs      int64
	EndToEndMs         int64
}

func RecordGenerationMetric(input GenerationMetricRecord) error {
	status := input.Status
	if input.Metrics != nil && input.Metrics.Status != "" {
		status = input.Metrics.Status
	}
	record := &db.KuroAIMetric{
		RequestID:          input.RequestID,
		ChannelID:          input.ChannelID,
		CreatedAt:          input.CreatedAt,
		Operation:          "generation",
		Status:             status,
		BotQueueMs:         input.BotQueueMs,
		DiscordHistoryMs:   input.DiscordHistoryMs,
		RuntimeRoundTripMs: input.RuntimeRoundTripMs,
		DiscordSendMs:      input.DiscordSendMs,
		EndToEndMs:         max(0, input.EndToEndMs),
	}
	applyGenerationMetrics(record, input.Metrics)
	return db.RecordKuroAIMetric(db.Dbs, record)
}

func RecordRuntimeMetric(event airuntime.MetricEvent) error {
	if event.RequestID == "" || event.ChannelID == "" || event.Metrics == nil {
		return nil
	}
	status := event.Metrics.Status
	if status == "" {
		status = "success"
	}
	operation := event.Operation
	if operation == "" {
		operation = "background"
	}
	record := &db.KuroAIMetric{
		RequestID: event.RequestID,
		ChannelID: event.ChannelID,
		CreatedAt: time.Now(),
		Operation: operation,
		Status:    status,
	}
	applyGenerationMetrics(record, event.Metrics)
	return db.RecordKuroAIMetric(db.Dbs, record)
}

func applyGenerationMetrics(record *db.KuroAIMetric, metrics *airuntime.GenerationMetrics) {
	if record == nil || metrics == nil {
		return
	}
	record.Model = metrics.Model
	record.Provider = metrics.Provider
	record.ProviderModel = metrics.ProviderModel
	record.RoutingStrategy = metrics.RoutingStrategy
	record.RoutingRegion = metrics.RoutingRegion
	record.RoutingAttempt = metrics.RoutingAttempt
	record.ProviderStatusCode = metrics.ProviderStatusCode
	record.UsageAvailable = metrics.UsageAvailable
	record.PromptTokens = metrics.PromptTokens
	record.CompletionTokens = metrics.CompletionTokens
	record.TotalTokens = metrics.TotalTokens
	record.ReasoningTokens = metrics.ReasoningTokens
	record.CachedTokens = metrics.CachedTokens
	record.CostUSD = metrics.CostUSD
	record.GenerationCount = metrics.GenerationCount
	record.RetryCount = metrics.RetryCount
	record.MemoryRecallMs = metrics.MemoryRecallMs
	record.BrowserQueueMs = metrics.BrowserQueueMs
	record.LocaleMs = metrics.LocaleMs
	record.PersonaMs = metrics.PersonaMs
	record.InjectUserMs = metrics.InjectUserMs
	record.GenerateMs = metrics.GenerateMs
	record.ValidateMs = metrics.ValidateMs
	record.RetryGenerateMs = metrics.RetryGenerateMs
	record.RetryValidateMs = metrics.RetryValidateMs
	record.BrowserTotalMs = metrics.BrowserTotalMs
	record.ProviderHeadersMs = metrics.ProviderHeadersMs
	record.ProviderFirstTokenMs = metrics.ProviderFirstTokenMs
	record.ProviderDurationMs = metrics.ProviderDurationMs
}
