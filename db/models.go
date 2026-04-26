package db

import "time"

type ZhtwToJp struct {
	ZhTw      string    `gorm:"primaryKey;size:1"` // 繁體中文漢字
	Jp        string    `gorm:"size:1;not null"`   // 日文漢字
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (ZhtwToJp) TableName() string {
	return "zhtw_to_jp"
}

// 誠也對應表，專門針對極端狀況去對應
type SeiyaCorrespond struct {
	GameName  string    `gorm:"primaryKey"`
	SeiyaURL  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type DiscordAllowList struct {
	ID         string    `gorm:"primaryKey"`
	Kind       string    `gorm:"not null"`
	Permission int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

type WebAPIToken struct {
	ID        string `gorm:"primaryKey"`
	ExpiresAt *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type (
	// game
	GameErogs struct {
		ID           int       `gorm:"primaryKey;autoIncrement:false" json:"id"`
		BrandErogsID int       `json:"brandErogsId"`
		Name         string    `gorm:"unique" json:"name"` // 遊戲名稱(批評空間)
		Image        string    `gorm:"not null;default:''" json:"image"`
		CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
		UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

		BrandErogs *BrandErogs `gorm:"foreignKey:BrandErogsID;references:ID" json:"brandErogs,omitempty"` // 單向 preload
	}
	// brand
	BrandErogs struct {
		ID        int       `gorm:"primaryKey;autoIncrement:false" json:"id"`
		Name      string    `gorm:"unique" json:"name"`
		Disband   bool      `json:"disband"`
		GameCount int       `gorm:"not null;default:0" json:"gameCount"`
		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
		UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	}
)

// 公告資料表
//
// 暫時沒有使用
type Announcement struct {
	ID        int       `gorm:"primaryKey"`
	Category  string    `gorm:"not null"`
	Content   string    `gorm:"not null"` // markdown內文
	Thumbnail *string   // 側邊小圖
	Image     *string   // 底部大圖
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

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
