package repository

import (
	"aniki/internal/database"

	"gorm.io/gorm"
)

type siteRepository struct {
	db *gorm.DB
}

// NewSiteRepository creates a new site repository instance
func NewSiteRepository(db *gorm.DB) SiteRepository {
	return &siteRepository{db: db}
}

func (r *siteRepository) Create(site *database.Site) error {
	return r.db.Create(site).Error
}

func (r *siteRepository) Update(site *database.Site) error {
	return r.db.Save(site).Error
}

func (r *siteRepository) FindByID(id int) (*database.Site, error) {
	var site database.Site
	err := r.db.First(&site, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &site, err
}

func (r *siteRepository) FindByName(name string) (*database.Site, error) {
	var site database.Site
	err := r.db.Where("name = ?", name).First(&site).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &site, err
}

func (r *siteRepository) FindAll() ([]database.Site, error) {
	var sites []database.Site
	err := r.db.Order("name").Find(&sites).Error
	return sites, err
}

func (r *siteRepository) Delete(id int) error {
	return r.db.Delete(&database.Site{}, id).Error
}
