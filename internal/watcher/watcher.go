package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"aniki/internal/database"
	"aniki/internal/hand_history"
	"aniki/internal/repository"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors directories for hand history files
type Watcher struct {
	watcher      *fsnotify.Watcher
	parser       *hand_history.Manager
	siteRepo     repository.SiteRepository
	handRepo     repository.HandRepository
	playerRepo   repository.PlayerRepository
	actionRepo   repository.ActionRepository
	paths        map[string]bool
	mu           sync.Mutex
	debounceMap  map[string]*time.Timer
	debounceMu   sync.Mutex
	stopCh       chan struct{}
	processingCh chan string
	workerCount  int
	isRunning    bool
}

// New creates a new file watcher
func New(parser *hand_history.Manager, siteRepo repository.SiteRepository, handRepo repository.HandRepository, playerRepo repository.PlayerRepository, actionRepo repository.ActionRepository) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		watcher:      fsWatcher,
		parser:       parser,
		siteRepo:     siteRepo,
		handRepo:     handRepo,
		playerRepo:   playerRepo,
		actionRepo:   actionRepo,
		paths:        make(map[string]bool),
		debounceMap:  make(map[string]*time.Timer),
		stopCh:       make(chan struct{}),
		processingCh: make(chan string, 100), // Buffer for file paths
		workerCount:  3,                      // 3 concurrent workers
		isRunning:    false,
	}, nil
}

// AddPath adds a directory to watch
func (w *Watcher) AddPath(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.paths[path] {
		return nil // Already watching
	}

	err := w.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("failed to watch path %s: %w", path, err)
	}

	w.paths[path] = true
	log.Printf("Now watching: %s", path)
	return nil
}

// RemovePath stops watching a directory
func (w *Watcher) RemovePath(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.paths[path] {
		return nil // Not watching
	}

	err := w.watcher.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to unwatch path %s: %w", path, err)
	}

	delete(w.paths, path)
	log.Printf("Stopped watching: %s", path)
	return nil
}

// Start begins watching for file changes
func (w *Watcher) Start() {
	w.mu.Lock()
	if w.isRunning {
		w.mu.Unlock()
		return
	}
	w.isRunning = true
	w.mu.Unlock()

	// Start worker pool
	for i := 0; i < w.workerCount; i++ {
		go w.worker(i)
	}

	// Start event listener
	go w.eventLoop()

	log.Println("File watcher started")
}

// Stop stops the file watcher
func (w *Watcher) Stop() {
	w.mu.Lock()
	if !w.isRunning {
		w.mu.Unlock()
		return
	}
	w.isRunning = false
	w.mu.Unlock()

	close(w.stopCh)
	w.watcher.Close()
	close(w.processingCh)

	log.Println("File watcher stopped")
}

// eventLoop listens for file system events
func (w *Watcher) eventLoop() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Only process write and create events
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				// Only process .txt files (PokerStars hand histories)
				if filepath.Ext(event.Name) == ".txt" {
					w.debounceFile(event.Name)
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)

		case <-w.stopCh:
			return
		}
	}
}

// debounceFile implements debouncing to avoid processing incomplete files
func (w *Watcher) debounceFile(filePath string) {
	w.debounceMu.Lock()
	defer w.debounceMu.Unlock()

	// Cancel existing timer for this file
	if timer, exists := w.debounceMap[filePath]; exists {
		timer.Stop()
	}

	// Create new timer that will trigger processing after delay
	w.debounceMap[filePath] = time.AfterFunc(1*time.Second, func() {
		w.processingCh <- filePath

		w.debounceMu.Lock()
		delete(w.debounceMap, filePath)
		w.debounceMu.Unlock()
	})
}

// worker processes files from the queue
func (w *Watcher) worker(id int) {
	for {
		select {
		case filePath, ok := <-w.processingCh:
			if !ok {
				return
			}
			w.processFile(filePath, id)

		case <-w.stopCh:
			return
		}
	}
}

// processFile parses and saves a hand history file
func (w *Watcher) processFile(filePath string, workerID int) {
	log.Printf("Worker %d: Processing file: %s", workerID, filePath)

	// Parse the file
	hands, siteName, err := w.parser.ParseFile(filePath)
	if err != nil {
		log.Printf("Worker %d: Error parsing file %s: %v", workerID, filePath, err)
		return
	}

	if len(hands) == 0 {
		log.Printf("Worker %d: No hands found in file: %s", workerID, filePath)
		return
	}

	// Get site from database
	site, err := w.siteRepo.FindByName(siteName)
	if err != nil {
		log.Printf("Worker %d: Error getting site %s: %v", workerID, siteName, err)
		return
	}
	if site == nil {
		log.Printf("Worker %d: Site not found: %s", workerID, siteName)
		return
	}

	// Process each hand
	saved := 0
	skipped := 0
	for _, hand := range hands {
		// Check if hand already exists
		exists, err := w.handRepo.Exists(site.ID, hand.HandID)
		if err != nil {
			log.Printf("Worker %d: Error checking hand existence: %v", workerID, err)
			continue
		}

		if exists {
			skipped++
			continue
		}

		// Convert hand_history.Hand to database.Hand
		dbHand := convertToDBHand(&hand, site.ID)

		// Save hand
		err = w.handRepo.Create(&dbHand)
		if err != nil {
			log.Printf("Worker %d: Error saving hand %s: %v", workerID, hand.HandID, err)
			continue
		}

		saved++
	}

	log.Printf("Worker %d: Processed %s - Saved: %d, Skipped: %d", workerID, filepath.Base(filePath), saved, skipped)
}

// convertToDBHand converts a hand_history.Hand to a database.Hand
func convertToDBHand(hand *hand_history.Hand, siteID int) database.Hand {
	// Convert hole cards to JSON
	holeCardsJSON, _ := json.Marshal(hand.HoleCards)

	// Convert board to JSON
	boardJSON, _ := json.Marshal(hand.Board)

	// Convert full hand to JSON for parsed_data
	parsedDataJSON, _ := json.Marshal(hand)

	return database.Hand{
		SiteID:     siteID,
		HandID:     hand.HandID,
		GameType:   hand.GameType,
		Stakes:     hand.Stakes,
		TableName:  hand.TableName,
		DateTime:   hand.DateTime,
		HeroName:   hand.HeroName,
		Position:   hand.Position,
		HoleCards:  string(holeCardsJSON),
		Board:      string(boardJSON),
		Result:     hand.Result,
		Rake:       hand.Rake,
		TotalPot:   hand.TotalPot,
		ParsedData: string(parsedDataJSON),
		RawText:    hand.RawText,
	}
}

// GetStatus returns the current status of the watcher
func (w *Watcher) GetStatus() map[string]interface{} {
	w.mu.Lock()
	defer w.mu.Unlock()

	return map[string]interface{}{
		"is_running":    w.isRunning,
		"watched_paths": w.paths,
		"queue_length":  len(w.processingCh),
	}
}
