package kuro

import (
	"testing"

	"kurohelperservice/airuntime"
	"kurohelperservice/db"
)

func TestApplyGenerationMetricsMapsRuntimeFields(t *testing.T) {
	record := &db.KuroAIMetric{}
	applyGenerationMetrics(record, &airuntime.GenerationMetrics{
		Model: "model-a", Provider: "provider-a", TotalTokens: 42,
		ReasoningTokens: 7, CostUSD: 0.001, ProviderFirstTokenMs: 900,
	})
	if record.Model != "model-a" || record.Provider != "provider-a" || record.TotalTokens != 42 {
		t.Fatalf("unexpected record: %#v", record)
	}
	if record.ReasoningTokens != 7 || record.CostUSD != 0.001 || record.ProviderFirstTokenMs != 900 {
		t.Fatalf("missing detailed metrics: %#v", record)
	}
}
