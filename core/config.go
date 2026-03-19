package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// DisplayMode represents the clock display mode.
type DisplayMode int

const (
	ModeDigital   DisplayMode = iota // "15:04:05"
	ModeAnalog                       // Analog clock face
	ModeTimestamp                    // Unix timestamp
)

// ThemeMode represents the application theme.
type ThemeMode int

const (
	ThemeLight ThemeMode = iota
	ThemeDark
	ThemeSystem
)

// WindowLevel represents the window z-order level.
type WindowLevel int

const (
	LevelTopMost WindowLevel = iota
	LevelNormal
)

// Config holds all persistent application settings.
type Config struct {
	DisplayMode  DisplayMode `json:"display_mode"`
	ThemeMode    ThemeMode   `json:"theme_mode"`
	WindowLevel  WindowLevel `json:"window_level"`
	PositionX    float32     `json:"position_x"`
	PositionY    float32     `json:"position_y"`
	Locked       bool        `json:"locked"`
	RestEnabled  bool        `json:"rest_enabled"`
	RestInterval int         `json:"rest_interval_minutes"`
	RestOpacity  float64     `json:"rest_opacity"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		DisplayMode:  ModeDigital,
		ThemeMode:    ThemeSystem,
		WindowLevel:  LevelTopMost,
		PositionX:    100,
		PositionY:    100,
		Locked:       false,
		RestEnabled:  true,
		RestInterval: 45,
		RestOpacity:  0.7,
	}
}

// configDir returns the platform-specific configuration directory.
func configDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Moment"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "Moment"), nil
	default:
		// Linux / other: use XDG_CONFIG_HOME or ~/.config
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			xdg = filepath.Join(home, ".config")
		}
		return filepath.Join(xdg, "Moment"), nil
	}
}

// ConfigStore manages loading and saving Config to a JSON file.
type ConfigStore struct {
	config   Config
	filePath string
	mu       sync.RWMutex
}

// NewConfigStore creates a ConfigStore with the platform-appropriate file path.
func NewConfigStore() (*ConfigStore, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	return &ConfigStore{
		config:   DefaultConfig(),
		filePath: filepath.Join(dir, "config.json"),
	}, nil
}

// NewConfigStoreWithPath creates a ConfigStore using a custom file path (useful for testing).
func NewConfigStoreWithPath(path string) *ConfigStore {
	return &ConfigStore{
		config:   DefaultConfig(),
		filePath: path,
	}
}

// Load reads the config from disk. If the file is missing or corrupted,
// it falls back to defaults and writes a fresh config file.
func (c *ConfigStore) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		// File missing — use defaults and persist them.
		c.config = DefaultConfig()
		return c.saveLocked()
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		// Corrupted — use defaults and overwrite.
		c.config = DefaultConfig()
		return c.saveLocked()
	}

	c.config = cfg
	return nil
}

// Save writes the current config to disk as JSON.
func (c *ConfigStore) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.saveLocked()
}

// saveLocked persists config without acquiring the lock (caller must hold it).
func (c *ConfigStore) saveLocked() error {
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c.config, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath, data, 0o644)
}

// Get returns a copy of the current config.
func (c *ConfigStore) Get() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// Update applies a mutation function to the config and saves it to disk.
func (c *ConfigStore) Update(fn func(*Config)) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fn(&c.config)
	return c.saveLocked()
}
