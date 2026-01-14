package repository

import (
	"time"

	"aniki/internal/database"
)

// SiteRepository defines the interface for site operations
type SiteRepository interface {
	Create(site *database.Site) error
	Update(site *database.Site) error
	FindByID(id int) (*database.Site, error)
	FindByName(name string) (*database.Site, error)
	FindAll() ([]database.Site, error)
	Delete(id int) error
}

// HandRepository defines the interface for hand operations
type HandRepository interface {
	Create(hand *database.Hand) error
	FindByID(id int64) (*database.Hand, error)
	FindAll(filter database.HandFilter) ([]database.Hand, error)
	Exists(siteID int, handID string) (bool, error)
	GetStats(heroName string) (*database.Stats, error)
	Delete(id int64) error
}

// PlayerRepository defines the interface for player operations
type PlayerRepository interface {
	Create(player *database.Player) error
	FindByHandID(handID int64) ([]database.Player, error)
	Delete(id int64) error
}

// ActionRepository defines the interface for action operations
type ActionRepository interface {
	Create(action *database.Action) error
	FindByHandID(handID int64) ([]database.Action, error)
	Delete(id int64) error
}

// HandFilter is used for querying hands
type HandFilter struct {
	SiteID   *int
	HeroName string
	GameType string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}
