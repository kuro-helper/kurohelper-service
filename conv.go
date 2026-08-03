package kurohelperservice

import (
	"kurohelperservice/db"
)

var (
	ZhtwToJp map[rune]rune
)

// init ZhtwToJp store
func InitZhtwToJp() error {
	entries, err := db.GetAllZhtwToJps(db.Dbs)
	if err != nil {
		return err
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
	return nil
}

func ZhTwToJp(search string) string {
	runes := []rune(search)
	for i, r := range runes {
		if jp, ok := ZhtwToJp[r]; ok {
			runes[i] = jp
		}
	}
	return string(runes)
}
