package database

import (
	"time"
)

// Site represents a poker site configuration
type Site struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"unique;not null"`
	WatchPath string    `json:"watch_path"`
	Enabled   bool      `json:"enabled" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Hands     []Hand    `json:"-" gorm:"foreignKey:SiteID;constraint:OnDelete:CASCADE"`
}

// Hand represents a parsed poker hand
type Hand struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SiteID     int       `json:"site_id" gorm:"not null;index"`
	Site       *Site     `json:"site,omitempty" gorm:"foreignKey:SiteID"`
	HandID     string    `json:"hand_id" gorm:"not null;uniqueIndex:idx_site_hand"`
	GameType   string    `json:"game_type" gorm:"index"`
	Stakes     string    `json:"stakes"`
	TableName  string    `json:"table_name"`
	DateTime   time.Time `json:"date_time" gorm:"index"`
	HeroName   string    `json:"hero_name" gorm:"index"`
	Position   string    `json:"position"`
	HoleCards  string    `json:"hole_cards"` // JSON array of cards
	Board      string    `json:"board"`      // JSON array of board cards
	Result     float64   `json:"result" gorm:"default:0"`
	Rake       float64   `json:"rake" gorm:"default:0"`
	TotalPot   float64   `json:"total_pot" gorm:"default:0"`
	ParsedData string    `json:"parsed_data" gorm:"type:text"`
	RawText    string    `json:"raw_text" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Players    []Player  `json:"players,omitempty" gorm:"foreignKey:HandID;constraint:OnDelete:CASCADE"`
	Actions    []Action  `json:"actions,omitempty" gorm:"foreignKey:HandID;constraint:OnDelete:CASCADE"`
}

// Player represents a player in a hand
type Player struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	HandID    int64     `json:"hand_id" gorm:"not null;index"`
	Hand      *Hand     `json:"-" gorm:"foreignKey:HandID"`
	Name      string    `json:"name" gorm:"not null"`
	Seat      int       `json:"seat"`
	Stack     float64   `json:"stack"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Action represents an action in a hand
type Action struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	HandID     int64     `json:"hand_id" gorm:"not null;index"`
	Hand       *Hand     `json:"-" gorm:"foreignKey:HandID"`
	PlayerName string    `json:"player_name" gorm:"not null"`
	Action     string    `json:"action" gorm:"not null"` // fold, check, call, bet, raise
	Amount     float64   `json:"amount" gorm:"default:0"`
	Street     string    `json:"street" gorm:"not null"` // preflop, flop, turn, river
	Sequence   int       `json:"sequence" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// HandFilter is used for querying hands
type HandFilter struct {
	SiteID   *int       `json:"site_id,omitempty"`
	HeroName string     `json:"hero_name,omitempty"`
	GameType string     `json:"game_type,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Limit    int        `json:"limit,omitempty"`
	Offset   int        `json:"offset,omitempty"`
}

// Stats represents aggregated statistics
type Stats struct {
	TotalHands  int     `json:"total_hands"`
	TotalWon    float64 `json:"total_won"`
	TotalRake   float64 `json:"total_rake"`
	BiggestWin  float64 `json:"biggest_win"`
	BiggestLoss float64 `json:"biggest_loss"`
	WinRate     float64 `json:"win_rate"` // BB/100
	HandsWon    int     `json:"hands_won"`
	HandsLost   int     `json:"hands_lost"`
}
