package db

import (
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const KuroAIMetricRetention = 30 * 24 * time.Hour

type KuroAIMetric struct {
	RequestID            string    `gorm:"primaryKey;column:request_id;size:64"`
	CreatedAt            time.Time `gorm:"index;not null"`
	ChannelID            string    `gorm:"index;size:32;not null"`
	Operation            string    `gorm:"index;size:32;not null;default:'generation'"`
	Status               string    `gorm:"index;size:16;not null"`
	Model                string    `gorm:"size:160"`
	Provider             string    `gorm:"index;size:120"`
	ProviderModel        string    `gorm:"size:200"`
	RoutingStrategy      string    `gorm:"size:32"`
	RoutingRegion        string    `gorm:"size:32"`
	RoutingAttempt       int
	ProviderStatusCode   int
	UsageAvailable       bool
	PromptTokens         int64
	CompletionTokens     int64
	TotalTokens          int64
	ReasoningTokens      int64
	CachedTokens         int64
	CostUSD              float64
	GenerationCount      int
	RetryCount           int
	BotQueueMs           int64
	DiscordHistoryMs     int64
	MemoryRecallMs       int64
	BrowserQueueMs       int64
	LocaleMs             int64
	PersonaMs            int64
	InjectUserMs         int64
	GenerateMs           int64
	ValidateMs           int64
	RetryGenerateMs      int64
	RetryValidateMs      int64
	BrowserTotalMs       int64
	ProviderHeadersMs    int64
	ProviderFirstTokenMs int64
	ProviderDurationMs   int64
	RuntimeRoundTripMs   int64
	DiscordSendMs        int64
	EndToEndMs           int64
}

type KuroAIDailyMetric struct {
	Day                   time.Time `gorm:"primaryKey;type:date"`
	RequestCount          int64
	SuccessCount          int64
	FailureCount          int64
	RetryCount            int64
	UsageCount            int64
	PromptTokens          int64
	CompletionTokens      int64
	TotalTokens           int64
	ReasoningTokens       int64
	CachedTokens          int64
	CostUSD               float64
	EndToEndMsTotal       int64
	RuntimeRoundTripTotal int64
	ProviderDurationTotal int64
	UpdatedAt             time.Time
}

func (KuroAIDailyMetric) TableName() string {
	return "kuro_ai_daily_metrics"
}

type KuroAIStats struct {
	RequestCount                     int64
	SuccessCount                     int64
	FailureCount                     int64
	RetryCount                       int64
	UsageCount                       int64
	PromptTokens                     int64
	CompletionTokens                 int64
	TotalTokens                      int64
	ReasoningTokens                  int64
	CachedTokens                     int64
	CostUSD                          float64
	AverageEndToEndMs                float64
	P50EndToEndMs                    float64
	P95EndToEndMs                    float64
	AverageRuntimeMs                 float64
	AverageProviderMs                float64
	MemoryExtractionCount            int64
	MemoryExtractionPromptTokens     int64
	MemoryExtractionCompletionTokens int64
	MemoryExtractionTotalTokens      int64
	MemoryExtractionCostUSD          float64
	VisionGenerationCount            int64
	VisionPromptTokens               int64
	VisionCompletionTokens           int64
	VisionTotalTokens                int64
	VisionCostUSD                    float64
}

type KuroAIProviderStats struct {
	Provider            string
	RequestCount        int64
	SuccessCount        int64
	FailureCount        int64
	AverageFirstTokenMs float64
	AverageDurationMs   float64
	P95DurationMs       float64
}

var lastKuroMetricCleanup atomic.Int64

func RecordKuroAIMetric(database *gorm.DB, metric *KuroAIMetric) error {
	if database == nil || metric == nil {
		return errors.New("DB or metric not initialized")
	}
	metric.RequestID = strings.TrimSpace(metric.RequestID)
	metric.ChannelID = strings.TrimSpace(metric.ChannelID)
	metric.Operation = strings.TrimSpace(metric.Operation)
	metric.Status = strings.TrimSpace(metric.Status)
	if metric.RequestID == "" || metric.ChannelID == "" || metric.Status == "" {
		return ErrParameterNotFound
	}
	if metric.CreatedAt.IsZero() {
		metric.CreatedAt = time.Now()
	}
	if metric.Operation == "" {
		metric.Operation = "generation"
	}

	return database.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(metric)
		if result.Error != nil || result.RowsAffected == 0 {
			return result.Error
		}

		day := time.Date(metric.CreatedAt.Year(), metric.CreatedAt.Month(), metric.CreatedAt.Day(), 0, 0, 0, 0, metric.CreatedAt.Location())
		success, failure, usage := int64(0), int64(0), int64(0)
		if metric.Status == "success" {
			success = 1
		} else {
			failure = 1
		}
		if metric.UsageAvailable {
			usage = 1
		}
		daily := KuroAIDailyMetric{
			Day: day, RequestCount: 1, SuccessCount: success, FailureCount: failure,
			RetryCount: int64(metric.RetryCount), UsageCount: usage,
			PromptTokens: metric.PromptTokens, CompletionTokens: metric.CompletionTokens,
			TotalTokens: metric.TotalTokens, ReasoningTokens: metric.ReasoningTokens,
			CachedTokens: metric.CachedTokens, CostUSD: metric.CostUSD,
			EndToEndMsTotal: metric.EndToEndMs, RuntimeRoundTripTotal: metric.RuntimeRoundTripMs,
			ProviderDurationTotal: metric.ProviderDurationMs, UpdatedAt: time.Now(),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "day"}},
			DoUpdates: clause.Assignments(map[string]any{
				"request_count":            gorm.Expr("kuro_ai_daily_metrics.request_count + EXCLUDED.request_count"),
				"success_count":            gorm.Expr("kuro_ai_daily_metrics.success_count + EXCLUDED.success_count"),
				"failure_count":            gorm.Expr("kuro_ai_daily_metrics.failure_count + EXCLUDED.failure_count"),
				"retry_count":              gorm.Expr("kuro_ai_daily_metrics.retry_count + EXCLUDED.retry_count"),
				"usage_count":              gorm.Expr("kuro_ai_daily_metrics.usage_count + EXCLUDED.usage_count"),
				"prompt_tokens":            gorm.Expr("kuro_ai_daily_metrics.prompt_tokens + EXCLUDED.prompt_tokens"),
				"completion_tokens":        gorm.Expr("kuro_ai_daily_metrics.completion_tokens + EXCLUDED.completion_tokens"),
				"total_tokens":             gorm.Expr("kuro_ai_daily_metrics.total_tokens + EXCLUDED.total_tokens"),
				"reasoning_tokens":         gorm.Expr("kuro_ai_daily_metrics.reasoning_tokens + EXCLUDED.reasoning_tokens"),
				"cached_tokens":            gorm.Expr("kuro_ai_daily_metrics.cached_tokens + EXCLUDED.cached_tokens"),
				"cost_usd":                 gorm.Expr("kuro_ai_daily_metrics.cost_usd + EXCLUDED.cost_usd"),
				"end_to_end_ms_total":      gorm.Expr("kuro_ai_daily_metrics.end_to_end_ms_total + EXCLUDED.end_to_end_ms_total"),
				"runtime_round_trip_total": gorm.Expr("kuro_ai_daily_metrics.runtime_round_trip_total + EXCLUDED.runtime_round_trip_total"),
				"provider_duration_total":  gorm.Expr("kuro_ai_daily_metrics.provider_duration_total + EXCLUDED.provider_duration_total"),
				"updated_at":               time.Now(),
			}),
		}).Create(&daily).Error; err != nil {
			return err
		}

		maybeCleanupKuroMetrics(tx, time.Now())
		return nil
	})
}

func maybeCleanupKuroMetrics(database *gorm.DB, now time.Time) {
	last := lastKuroMetricCleanup.Load()
	if last != 0 && now.Unix()-last < int64(time.Hour/time.Second) {
		return
	}
	if !lastKuroMetricCleanup.CompareAndSwap(last, now.Unix()) {
		return
	}
	_ = database.Where("created_at < ?", now.Add(-KuroAIMetricRetention)).Delete(&KuroAIMetric{}).Error
}

func GetKuroAIStats(database *gorm.DB, since time.Time) (KuroAIStats, error) {
	var stats KuroAIStats
	if database == nil {
		return stats, errors.New("DB not initialized")
	}
	err := database.Raw(`
		SELECT
			COUNT(*) AS request_count,
			COUNT(*) FILTER (WHERE status = 'success') AS success_count,
			COUNT(*) FILTER (WHERE status <> 'success') AS failure_count,
			COALESCE(SUM(retry_count), 0) AS retry_count,
			COUNT(*) FILTER (WHERE usage_available) AS usage_count,
			COALESCE(SUM(prompt_tokens), 0) AS prompt_tokens,
			COALESCE(SUM(completion_tokens), 0) AS completion_tokens,
			COALESCE(SUM(total_tokens), 0) AS total_tokens,
			COALESCE(SUM(reasoning_tokens), 0) AS reasoning_tokens,
			COALESCE(SUM(cached_tokens), 0) AS cached_tokens,
			COALESCE(SUM(cost_usd), 0) AS cost_usd,
			COALESCE(AVG(end_to_end_ms) FILTER (WHERE status = 'success' AND operation = 'generation'), 0) AS average_end_to_end_ms,
			COALESCE(PERCENTILE_CONT(0.50) WITHIN GROUP (ORDER BY end_to_end_ms) FILTER (WHERE status = 'success' AND operation = 'generation'), 0) AS p50_end_to_end_ms,
			COALESCE(PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY end_to_end_ms) FILTER (WHERE status = 'success' AND operation = 'generation'), 0) AS p95_end_to_end_ms,
			COALESCE(AVG(runtime_round_trip_ms) FILTER (WHERE status = 'success' AND operation = 'generation'), 0) AS average_runtime_ms,
			COALESCE(AVG(provider_duration_ms) FILTER (WHERE status = 'success' AND operation = 'generation' AND provider_duration_ms > 0), 0) AS average_provider_ms,
			COUNT(*) FILTER (WHERE operation = 'memory_extraction') AS memory_extraction_count,
			COALESCE(SUM(prompt_tokens) FILTER (WHERE operation = 'memory_extraction'), 0) AS memory_extraction_prompt_tokens,
			COALESCE(SUM(completion_tokens) FILTER (WHERE operation = 'memory_extraction'), 0) AS memory_extraction_completion_tokens,
			COALESCE(SUM(total_tokens) FILTER (WHERE operation = 'memory_extraction'), 0) AS memory_extraction_total_tokens,
			COALESCE(SUM(cost_usd) FILTER (WHERE operation = 'memory_extraction'), 0) AS memory_extraction_cost_usd,
			COALESCE(SUM(GREATEST(generation_count, 1)) FILTER (WHERE operation = 'vision'), 0) AS vision_generation_count,
			COALESCE(SUM(prompt_tokens) FILTER (WHERE operation = 'vision'), 0) AS vision_prompt_tokens,
			COALESCE(SUM(completion_tokens) FILTER (WHERE operation = 'vision'), 0) AS vision_completion_tokens,
			COALESCE(SUM(total_tokens) FILTER (WHERE operation = 'vision'), 0) AS vision_total_tokens,
			COALESCE(SUM(cost_usd) FILTER (WHERE operation = 'vision'), 0) AS vision_cost_usd
		FROM kuro_ai_metrics
		WHERE created_at >= ?
	`, since).Scan(&stats).Error
	return stats, err
}

func GetKuroAIProviderStats(database *gorm.DB, since time.Time) ([]KuroAIProviderStats, error) {
	if database == nil {
		return nil, errors.New("DB not initialized")
	}
	var providers []KuroAIProviderStats
	err := database.Raw(`
		SELECT
			COALESCE(NULLIF(provider, ''), 'unknown') AS provider,
			COUNT(*) AS request_count,
			COUNT(*) FILTER (WHERE status = 'success') AS success_count,
			COUNT(*) FILTER (WHERE status <> 'success') AS failure_count,
			COALESCE(AVG(provider_first_token_ms) FILTER (WHERE provider_first_token_ms > 0), 0) AS average_first_token_ms,
			COALESCE(AVG(provider_duration_ms) FILTER (WHERE provider_duration_ms > 0), 0) AS average_duration_ms,
			COALESCE(PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY provider_duration_ms) FILTER (WHERE provider_duration_ms > 0), 0) AS p95_duration_ms
		FROM kuro_ai_metrics
		WHERE created_at >= ?
		GROUP BY COALESCE(NULLIF(provider, ''), 'unknown')
		ORDER BY request_count DESC, provider ASC
	`, since).Scan(&providers).Error
	return providers, err
}
