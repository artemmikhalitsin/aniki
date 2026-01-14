package repository

import (
	"aniki/internal/database"

	"gorm.io/gorm"
)

type handRepository struct {
	db *gorm.DB
}

// NewHandRepository creates a new hand repository instance
func NewHandRepository(db *gorm.DB) HandRepository {
	return &handRepository{db: db}
}

func (r *handRepository) Create(hand *database.Hand) error {
	return r.db.Create(hand).Error
}

func (r *handRepository) FindByID(id int64) (*database.Hand, error) {
	var hand database.Hand
	err := r.db.Preload("Site").Preload("Players").Preload("Actions").First(&hand, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &hand, err
}

func (r *handRepository) FindAll(filter database.HandFilter) ([]database.Hand, error) {
	var hands []database.Hand
	query := r.db.Model(&database.Hand{})

	if filter.SiteID != nil {
		query = query.Where("site_id = ?", *filter.SiteID)
	}

	if filter.HeroName != "" {
		query = query.Where("hero_name = ?", filter.HeroName)
	}

	if filter.GameType != "" {
		query = query.Where("game_type = ?", filter.GameType)
	}

	if filter.DateFrom != nil {
		query = query.Where("date_time >= ?", *filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("date_time <= ?", *filter.DateTo)
	}

	query = query.Order("date_time DESC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&hands).Error
	return hands, err
}

func (r *handRepository) Exists(siteID int, handID string) (bool, error) {
	var count int64
	err := r.db.Model(&database.Hand{}).Where("site_id = ? AND hand_id = ?", siteID, handID).Count(&count).Error
	return count > 0, err
}

func (r *handRepository) GetStats(heroName string) (*database.Stats, error) {
	stats := &database.Stats{}

	type Result struct {
		TotalHands  int64
		TotalWon    float64
		TotalRake   float64
		BiggestWin  float64
		BiggestLoss float64
		HandsWon    int64
		HandsLost   int64
	}

	var result Result
	err := r.db.Model(&database.Hand{}).
		Where("hero_name = ?", heroName).
		Select(`
			COUNT(*) as total_hands,
			COALESCE(SUM(result), 0) as total_won,
			COALESCE(SUM(rake), 0) as total_rake,
			COALESCE(MAX(result), 0) as biggest_win,
			COALESCE(MIN(result), 0) as biggest_loss,
			COALESCE(SUM(CASE WHEN result > 0 THEN 1 ELSE 0 END), 0) as hands_won,
			COALESCE(SUM(CASE WHEN result < 0 THEN 1 ELSE 0 END), 0) as hands_lost
		`).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	stats.TotalHands = int(result.TotalHands)
	stats.TotalWon = result.TotalWon
	stats.TotalRake = result.TotalRake
	stats.BiggestWin = result.BiggestWin
	stats.BiggestLoss = result.BiggestLoss
	stats.HandsWon = int(result.HandsWon)
	stats.HandsLost = int(result.HandsLost)

	// Calculate win rate (simplified - would need more context for accurate BB/100)
	if stats.TotalHands > 0 {
		stats.WinRate = (stats.TotalWon / float64(stats.TotalHands)) * 100
	}

	return stats, nil
}

func (r *handRepository) Delete(id int64) error {
	return r.db.Delete(&database.Hand{}, id).Error
}
