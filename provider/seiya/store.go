package seiya

import (
	"kurohelperservice/db"
	"log/slog"
	"os"
)

var (
	SeiyaCorrespond map[string]string
)

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
