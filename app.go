package main

import (
	"context"
	"fmt"
	"log"

	"aniki/internal/config"
	"aniki/internal/database"
	"aniki/internal/hand_history"
	"aniki/internal/repository"
	"aniki/internal/watcher"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	db         *database.DB
	config     *config.Config
	watcher    *watcher.Watcher
	parser     *hand_history.Manager
	siteRepo   repository.SiteRepository
	handRepo   repository.HandRepository
	playerRepo repository.PlayerRepository
	actionRepo repository.ActionRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		cfg, _ = config.GetDefaultConfig()
	}
	a.config = cfg

	// Initialize database
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		log.Fatalf("Failed to get database path: %v", err)
	}

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	a.db = db

	// Initialize repositories
	a.siteRepo = repository.NewSiteRepository(db.DB)
	a.handRepo = repository.NewHandRepository(db.DB)
	a.playerRepo = repository.NewPlayerRepository(db.DB)
	a.actionRepo = repository.NewActionRepository(db.DB)

	// Initialize parser
	a.parser = hand_history.NewManager()

	// Initialize file watcher
	w, err := watcher.New(a.parser, a.siteRepo, a.handRepo, a.playerRepo, a.actionRepo)
	if err != nil {
		log.Fatalf("Failed to initialize watcher: %v", err)
	}
	a.watcher = w

	// Initialize default sites in database
	a.initializeSites()

	// Start watching configured paths
	for _, site := range a.config.Sites {
		if site.Enabled && site.WatchPath != "" {
			if err := a.watcher.AddPath(site.WatchPath); err != nil {
				log.Printf("Failed to watch path %s: %v", site.WatchPath, err)
			}
		}
	}

	a.watcher.Start()

	log.Println("Application started successfully")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.watcher != nil {
		a.watcher.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
	log.Println("Application shutdown complete")
}

// initializeSites ensures default sites exist in the database
func (a *App) initializeSites() {
	for _, siteCfg := range a.config.Sites {
		site, err := a.siteRepo.FindByName(siteCfg.Name)
		if err != nil {
			log.Printf("Error checking site %s: %v", siteCfg.Name, err)
			continue
		}

		if site == nil {
			// Site doesn't exist, create it
			newSite := &database.Site{
				Name:      siteCfg.Name,
				WatchPath: siteCfg.WatchPath,
				Enabled:   siteCfg.Enabled,
			}
			if err := a.siteRepo.Create(newSite); err != nil {
				log.Printf("Error creating site %s: %v", siteCfg.Name, err)
			} else {
				log.Printf("Created site: %s", siteCfg.Name)
			}
		}
	}
}

// GetHands retrieves hands based on the provided filter
func (a *App) GetHands(filter database.HandFilter) ([]database.Hand, error) {
	return a.handRepo.FindAll(filter)
}

// GetHandByID retrieves a single hand by ID
func (a *App) GetHandByID(id int64) (*database.Hand, error) {
	return a.handRepo.FindByID(id)
}

// GetStats retrieves statistics for a hero
func (a *App) GetStats(heroName string) (*database.Stats, error) {
	return a.handRepo.GetStats(heroName)
}

// GetConfig returns the current configuration
func (a *App) GetConfig() *config.Config {
	return a.config
}

// UpdateConfig updates the configuration
func (a *App) UpdateConfig(cfg *config.Config) error {
	// Stop watching old paths
	for _, site := range a.config.Sites {
		if site.WatchPath != "" {
			a.watcher.RemovePath(site.WatchPath)
		}
	}

	// Save new config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	a.config = cfg

	// Update database sites
	for _, siteCfg := range cfg.Sites {
		site, err := a.siteRepo.FindByName(siteCfg.Name)
		if err != nil {
			log.Printf("Error getting site %s: %v", siteCfg.Name, err)
			continue
		}

		if site != nil {
			site.WatchPath = siteCfg.WatchPath
			site.Enabled = siteCfg.Enabled
			if err := a.siteRepo.Update(site); err != nil {
				log.Printf("Error updating site %s: %v", siteCfg.Name, err)
			}
		}
	}

	// Start watching new paths
	for _, site := range cfg.Sites {
		if site.Enabled && site.WatchPath != "" {
			if err := a.watcher.AddPath(site.WatchPath); err != nil {
				log.Printf("Failed to watch path %s: %v", site.WatchPath, err)
			}
		}
	}

	return nil
}

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Hand History Directory",
	})
	return path, err
}

// GetSites retrieves all sites from the database
func (a *App) GetSites() ([]database.Site, error) {
	return a.siteRepo.FindAll()
}

// GetWatcherStatus returns the current watcher status
func (a *App) GetWatcherStatus() map[string]interface{} {
	return a.watcher.GetStatus()
}
