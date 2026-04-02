package store

import (
	"log/slog"
	"os"

	"kurohelperservice/db"
)

var (
	ZhtwToJp        map[rune]rune
	SeiyaCorrespond map[string]string
)

// init ZhtwToJp store
func InitZhtwToJp() {
	entries, err := db.GetAllZhtwToJps(db.Dbs)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// 轉換
	ZhtwToJp = make(map[rune]rune, len(entries))
	for _, e := range entries {
		keyRunes := []rune(e.ZhTw)
		valRunes := []rune(e.Jp)

		// 確保都是單一字
		if len(keyRunes) == 1 && len(valRunes) == 1 {
			ZhtwToJp[keyRunes[0]] = valRunes[0]
		}
	}
}

// init SeiyaCorrespond store
func InitSeiyaCorrespond() {
	entries, err := db.GetAllSeiyaCorresponds(db.Dbs)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Translate
	SeiyaCorrespond = make(map[string]string, len(entries))
	for _, e := range entries {
		SeiyaCorrespond[e.GameName] = e.SeiyaURL
	}
}
