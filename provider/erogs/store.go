package erogs

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"

	"kurohelperservice"
)

var (
	GamesName     []string       = make([]string, 0)
	InvertedIndex map[rune][]int = make(map[rune][]int)
)

func InitErogsGameAutoComplete(pathRoute string) {
	jsonFile, err := os.Open(pathRoute)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var entries []GameItem

	err = json.Unmarshal(byteValue, &entries)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	GamesName = make([]string, 0, len(entries))
	for _, e := range entries {
		GamesName = append(GamesName, e.Name)
	}

	GamesName = kurohelperservice.Distinct(GamesName)

	// 初始化倒排索引
	InvertedIndex = make(map[rune][]int)

	for i, name := range GamesName {
		runes := []rune(strings.ToLower(name))

		// 避免重複存入索引
		seen := make(map[rune]bool)

		for _, char := range runes {
			if !seen[char] {
				InvertedIndex[char] = append(InvertedIndex[char], i)
				seen[char] = true
			}
		}
	}
}
