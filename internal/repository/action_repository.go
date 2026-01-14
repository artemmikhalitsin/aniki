package repository

import (
	"aniki/internal/database"

	"gorm.io/gorm"
)

type actionRepository struct {
	db *gorm.DB
}

// NewActionRepository creates a new action repository instance
func NewActionRepository(db *gorm.DB) ActionRepository {
	return &actionRepository{db: db}
}

func (r *actionRepository) Create(action *database.Action) error {
	return r.db.Create(action).Error
}

func (r *actionRepository) FindByHandID(handID int64) ([]database.Action, error) {
	var actions []database.Action
	err := r.db.Where("hand_id = ?", handID).Order("sequence").Find(&actions).Error
	return actions, err
}

func (r *actionRepository) Delete(id int64) error {
	return r.db.Delete(&database.Action{}, id).Error
}
