package seiya

import (
	"kurohelperservice/db"
)

var (
	SeiyaCorrespond map[string]string
)

// init SeiyaCorrespond store
func InitSeiyaCorrespond() error {
	entries, err := db.GetAllSeiyaCorresponds(db.Dbs)
	if err != nil {
		return err
	}

	// Translate
	SeiyaCorrespond = make(map[string]string, len(entries))
	for _, e := range entries {
		SeiyaCorrespond[e.GameName] = e.SeiyaURL
	}
	return nil
}
