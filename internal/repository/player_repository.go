package repository

import (
	"aniki/internal/database"

	"gorm.io/gorm"
)

type playerRepository struct {
	db *gorm.DB
}

// NewPlayerRepository creates a new player repository instance
func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(player *database.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) FindByHandID(handID int64) ([]database.Player, error) {
	var players []database.Player
	err := r.db.Where("hand_id = ?", handID).Find(&players).Error
	return players, err
}

func (r *playerRepository) Delete(id int64) error {
	return r.db.Delete(&database.Player{}, id).Error
}
