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
	GamesName          []string       = make([]string, 0)
	GameInvertedIndex  map[rune][]int = make(map[rune][]int)
	BrandsName         []string       = make([]string, 0)
	BrandInvertedIndex map[rune][]int = make(map[rune][]int)
	MusicsName         []string       = make([]string, 0)
	MusicInvertedIndex map[rune][]int = make(map[rune][]int)
)

func loadErogsAutocompleteFromJSON(path string) ([]string, map[rune][]int, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, nil, err
	}

	var entries []GameItem
	if err := json.Unmarshal(byteValue, &entries); err != nil {
		return nil, nil, err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name)
	}
	names = kurohelperservice.Distinct(names)

	idx := make(map[rune][]int)
	for i, name := range names {
		runes := []rune(strings.ToLower(name))
		seen := make(map[rune]bool)
		for _, char := range runes {
			if !seen[char] {
				idx[char] = append(idx[char], i)
				seen[char] = true
			}
		}
	}
	return names, idx, nil
}

func InitErogsGameAutoComplete(pathRoute string) {
	names, idx, err := loadErogsAutocompleteFromJSON(pathRoute)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	GamesName = names
	GameInvertedIndex = idx
}

func InitErogsBrandAutoComplete(pathRoute string) {
	names, idx, err := loadErogsAutocompleteFromJSON(pathRoute)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	BrandsName = names
	BrandInvertedIndex = idx
}

func InitErogsMusicAutoComplete(pathRoute string) {
	names, idx, err := loadErogsAutocompleteFromJSON(pathRoute)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	MusicsName = names
	MusicInvertedIndex = idx
}
