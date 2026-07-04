package db

import "time"

// 未來用
type Game struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	PlayHours   float64
	Brand       int // preload
	Links       int // preload
	Description string
	ScoreAvg    float64
	ScoreCount  int
	ReleaseDate time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	CreatedUser int
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	UpdatedUser int
	BangumiID   int
	VNDBID      int
	ErogsID     int
	Alias       int // preload
	Tag         int // preload
}
