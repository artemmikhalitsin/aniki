package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds the application configuration
type Config struct {
	HeroName     string          `json:"hero_name"`
	Sites        map[string]Site `json:"sites"`
	DatabasePath string          `json:"database_path"`
	Theme        string          `json:"theme"`
}

// Site represents per-site configuration
type Site struct {
	Name      string `json:"name"`
	WatchPath string `json:"watch_path"`
	Enabled   bool   `json:"enabled"`
}

// GetConfigDir returns the platform-specific configuration directory
func GetConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "Aniki")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(home, "Library", "Application Support", "Aniki")
	case "linux":
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			xdgConfig = filepath.Join(home, ".config")
		}
		configDir = filepath.Join(xdgConfig, "aniki")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// LoadConfig loads the configuration from disk
func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.json")

	// Return default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	defaultWatchPath := DetectPokerStarsPath()

	return &Config{
		HeroName: "",
		Sites: map[string]Site{
			"pokerstars": {
				Name:      "PokerStars",
				WatchPath: defaultWatchPath,
				Enabled:   true,
			},
		},
		DatabasePath: filepath.Join(configDir, "poker.db"),
		Theme:        "dark",
	}, nil
}

// DetectPokerStarsPath attempts to detect the PokerStars hand history directory
func DetectPokerStarsPath() string {
	var basePath string

	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		basePath = filepath.Join(localAppData, "PokerStars", "HandHistory")
	case "darwin":
		home, _ := os.UserHomeDir()
		basePath = filepath.Join(home, "Library", "Application Support", "PokerStars", "HandHistory")
	case "linux":
		home, _ := os.UserHomeDir()
		// Try Wine path first
		winePath := filepath.Join(home, ".wine", "drive_c", "users", os.Getenv("USER"),
			"Local Settings", "Application Data", "PokerStars", "HandHistory")
		if _, err := os.Stat(winePath); err == nil {
			return winePath
		}
		// Try native path
		basePath = filepath.Join(home, ".local", "share", "PokerStars", "HandHistory")
	}

	// Check if path exists
	if _, err := os.Stat(basePath); err == nil {
		return basePath
	}

	return ""
}

// GetDatabasePath returns the full path to the database file
func GetDatabasePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "poker.db"), nil
}
