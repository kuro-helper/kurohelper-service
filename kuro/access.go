package kuro

import "strings"

type TextCommand struct {
	Name string
	Args []string
}

func ParseIDSet(value string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result[item] = struct{}{}
		}
	}
	return result
}

func IsAllowed(ids map[string]struct{}, id string) bool {
	if len(ids) == 0 {
		return true
	}
	_, ok := ids[id]
	return ok
}

func PrepareTrigger(content, prefix, botUserID string, mentioned bool) (string, bool) {
	content = strings.TrimSpace(content)
	hasPrefix := prefix != "" && strings.HasPrefix(content, prefix)
	if prefix != "" && !hasPrefix && !mentioned {
		return content, false
	}
	if mentioned && botUserID != "" {
		content = strings.ReplaceAll(content, "<@"+botUserID+">", " ")
		content = strings.ReplaceAll(content, "<@!"+botUserID+">", " ")
		content = strings.Join(strings.Fields(content), " ")
	}
	return strings.TrimSpace(content), true
}

// ParseTextCommand parses Kuro management commands in the form
// "<trigger prefix> /<command> [arguments]". It intentionally does not parse
// bare slash commands or bot mentions so management commands have one
// unambiguous Discord text format.
func ParseTextCommand(content, prefix string) (TextCommand, bool) {
	content = strings.TrimSpace(content)
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || !strings.HasPrefix(content, prefix) {
		return TextCommand{}, false
	}

	remainder := strings.TrimPrefix(content, prefix)
	if remainder != "" {
		first := []rune(remainder)[0]
		if first != '/' && first != ' ' && first != '\t' && first != '\r' && first != '\n' {
			return TextCommand{}, false
		}
	}
	remainder = strings.TrimSpace(remainder)
	if !strings.HasPrefix(remainder, "/") {
		return TextCommand{}, false
	}

	fields := strings.Fields(strings.TrimPrefix(remainder, "/"))
	if len(fields) == 0 {
		return TextCommand{Name: "help"}, true
	}

	name := strings.ToLower(strings.TrimSpace(fields[0]))
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}
	return TextCommand{Name: name, Args: args}, true
}
