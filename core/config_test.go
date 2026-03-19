package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempConfigPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "config.json")
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DisplayMode != ModeDigital {
		t.Errorf("expected ModeDigital, got %d", cfg.DisplayMode)
	}
	if cfg.ThemeMode != ThemeSystem {
		t.Errorf("expected ThemeSystem, got %d", cfg.ThemeMode)
	}
	if cfg.RestInterval != 45 {
		t.Errorf("expected 45, got %d", cfg.RestInterval)
	}
	if cfg.RestOpacity != 0.7 {
		t.Errorf("expected 0.7, got %f", cfg.RestOpacity)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempConfigPath(t)
	store := NewConfigStoreWithPath(path)

	// Update a field and save
	err := store.Update(func(c *Config) {
		c.DisplayMode = ModeAnalog
		c.ThemeMode = ThemeDark
		c.PositionX = 200
		c.PositionY = 300
	})
	if err != nil {
		t.Fatal(err)
	}

	// Load into a new store
	store2 := NewConfigStoreWithPath(path)
	if err := store2.Load(); err != nil {
		t.Fatal(err)
	}

	got := store2.Get()
	if got.DisplayMode != ModeAnalog {
		t.Errorf("expected ModeAnalog, got %d", got.DisplayMode)
	}
	if got.ThemeMode != ThemeDark {
		t.Errorf("expected ThemeDark, got %d", got.ThemeMode)
	}
	if got.PositionX != 200 || got.PositionY != 300 {
		t.Errorf("expected (200,300), got (%f,%f)", got.PositionX, got.PositionY)
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent", "config.json")
	store := NewConfigStoreWithPath(path)

	if err := store.Load(); err != nil {
		t.Fatal(err)
	}

	// Should have defaults
	got := store.Get()
	expected := DefaultConfig()
	if got != expected {
		t.Errorf("expected defaults, got %+v", got)
	}

	// File should now exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestLoadCorruptedFile(t *testing.T) {
	path := tempConfigPath(t)
	// Write garbage
	if err := os.WriteFile(path, []byte("{invalid json!!!"), 0o644); err != nil {
		t.Fatal(err)
	}

	store := NewConfigStoreWithPath(path)
	if err := store.Load(); err != nil {
		t.Fatal(err)
	}

	got := store.Get()
	expected := DefaultConfig()
	if got != expected {
		t.Errorf("expected defaults after corrupt load, got %+v", got)
	}

	// Verify the file was overwritten with valid JSON
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var check Config
	if err := json.Unmarshal(data, &check); err != nil {
		t.Errorf("config file should contain valid JSON after recovery: %v", err)
	}
}

func TestGetReturnsCopy(t *testing.T) {
	store := NewConfigStoreWithPath(tempConfigPath(t))
	cfg1 := store.Get()
	cfg1.PositionX = 999
	cfg2 := store.Get()
	if cfg2.PositionX == 999 {
		t.Error("Get() should return a copy, not a reference")
	}
}
