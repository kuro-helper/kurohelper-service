package kuro

import (
	"strings"
	"testing"
	"time"
)

func TestBuildRecentContextKeepsSpeakersAndBoundary(t *testing.T) {
	now := time.Now()
	prompt, retrieval := BuildRecentContext([]RecentMessage{
		{ID: "100", DisplayName: "舊使用者", Content: "不應出現", CreatedAt: now},
		{ID: "102", DisplayName: "肉圓", Content: "早安", CreatedAt: now.Add(time.Second)},
		{ID: "103", DisplayName: "Kuro", Content: "……早安。", Assistant: true, CreatedAt: now.Add(2 * time.Second)},
	}, ContextOptions{BoundaryID: "101", MessageLimit: 15, MaxChars: 6000})
	if strings.Contains(prompt, "不應出現") {
		t.Fatal("context included a message before the new-chat boundary")
	}
	if !strings.Contains(retrieval, "[肉圓] 早安") || !strings.Contains(retrieval, "[Kuro] ……早安。") {
		t.Fatalf("unexpected context: %s", retrieval)
	}
}

func TestPrepareTriggerSupportsPrefixAndMention(t *testing.T) {
	if text, ok := PrepareTrigger("小黑早安", "小黑", "bot", false); !ok || text != "小黑早安" {
		t.Fatalf("prefix was not preserved: %q %v", text, ok)
	}
	if text, ok := PrepareTrigger("<@bot> 早安", "小黑", "bot", true); !ok || text != "早安" {
		t.Fatalf("mention was not removed: %q %v", text, ok)
	}
	if _, ok := PrepareTrigger("大家早安", "小黑", "bot", false); ok {
		t.Fatal("unaddressed message should not trigger")
	}
}

func TestCollectContextParticipantsUsesStableIDsAndDeduplicates(t *testing.T) {
	participants := CollectContextParticipants(
		[]RecentMessage{
			{UserID: "u2", DisplayName: "海獺"},
			{UserID: "bot", DisplayName: "Kuro", Assistant: true},
			{UserID: "u1", DisplayName: "肉圓舊名"},
		},
		MentionedUser{ID: "u1", DisplayName: "肉圓"},
		[]MentionedUser{{ID: "u3", DisplayName: "Tommy"}},
	)
	if len(participants) != 3 {
		t.Fatalf("participants = %#v", participants)
	}
	if participants[0].ID != "u1" || participants[0].DisplayName != "肉圓" {
		t.Fatalf("current speaker did not keep priority: %#v", participants[0])
	}
	if participants[1].ID != "u3" || participants[2].ID != "u2" {
		t.Fatalf("unexpected participant order: %#v", participants)
	}
}
