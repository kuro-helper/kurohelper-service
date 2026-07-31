package kuro

import (
	"encoding/json"
	"testing"
)

func TestGenerateResponseDecodesMetrics(t *testing.T) {
	var response GenerateResponse
	err := json.Unmarshal([]byte(`{
		"text":"reply",
		"metrics":{"status":"success","model":"model-a","provider":"Provider B","providerModel":"native-b","routingStrategy":"fallback","routingAttempt":2,"totalTokens":42,"costUsd":0.001,"providerDurationMs":900}
	}`), &response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Metrics == nil || response.Metrics.TotalTokens != 42 || response.Metrics.CostUSD != 0.001 {
		t.Fatalf("unexpected metrics: %#v", response.Metrics)
	}
	if response.Metrics.Provider != "Provider B" || response.Metrics.RoutingAttempt != 2 {
		t.Fatalf("unexpected routing metrics: %#v", response.Metrics)
	}
}
