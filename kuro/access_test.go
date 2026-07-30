package kuro

import (
	"reflect"
	"testing"
)

func TestParseTextCommand(t *testing.T) {
	tests := []struct {
		name    string
		content string
		prefix  string
		want    TextCommand
		ok      bool
	}{
		{name: "English new chat", content: "小黑 /newchat", prefix: "小黑", want: TextCommand{Name: "newchat"}, ok: true},
		{name: "Explicit help", content: "小黑 /help", prefix: "小黑", want: TextCommand{Name: "help"}, ok: true},
		{name: "Memory list page", content: "小黑 /memory-list 30", prefix: "小黑", want: TextCommand{Name: "memory-list", Args: []string{"30"}}, ok: true},
		{name: "Forget memory", content: "小黑 /forget abcdef12", prefix: "小黑", want: TextCommand{Name: "forget", Args: []string{"abcdef12"}}, ok: true},
		{name: "Chinese command is not an alias", content: "小黑 /新對話", prefix: "小黑", want: TextCommand{Name: "新對話"}, ok: true},
		{name: "Compact form", content: "小黑/status", prefix: "小黑", want: TextCommand{Name: "status"}, ok: true},
		{name: "Empty command is help", content: "小黑 /", prefix: "小黑", want: TextCommand{Name: "help"}, ok: true},
		{name: "Unknown command remains parseable", content: "小黑 /unknown", prefix: "小黑", want: TextCommand{Name: "unknown"}, ok: true},
		{name: "Ordinary chat", content: "小黑 今天好嗎", prefix: "小黑", ok: false},
		{name: "Bare slash command", content: "/newchat", prefix: "小黑", ok: false},
		{name: "Prefix collision", content: "小黑貓 /newchat", prefix: "小黑", ok: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := ParseTextCommand(test.content, test.prefix)
			if ok != test.ok || !reflect.DeepEqual(got, test.want) {
				t.Fatalf("ParseTextCommand() = (%#v, %t), want (%#v, %t)", got, ok, test.want, test.ok)
			}
		})
	}
}
