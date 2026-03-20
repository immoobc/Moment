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
	if cfg.WindowLevel != LevelTopMost {
		t.Errorf("expected LevelTopMost, got %d", cfg.WindowLevel)
	}
	if cfg.Locked {
		t.Error("expected Locked=false")
	}
	if cfg.Theme != ThemeLight {
		t.Errorf("expected ThemeLight, got %d", cfg.Theme)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempConfigPath(t)
	store := NewConfigStoreWithPath(path)

	err := store.Update(func(c *Config) {
		c.PositionX = 200
		c.PositionY = 300
		c.Locked = true
	})
	if err != nil {
		t.Fatal(err)
	}

	store2 := NewConfigStoreWithPath(path)
	if err := store2.Load(); err != nil {
		t.Fatal(err)
	}

	got := store2.Get()
	if got.PositionX != 200 || got.PositionY != 300 {
		t.Errorf("expected (200,300), got (%f,%f)", got.PositionX, got.PositionY)
	}
	if !got.Locked {
		t.Error("expected Locked=true")
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent", "config.json")
	store := NewConfigStoreWithPath(path)

	if err := store.Load(); err != nil {
		t.Fatal(err)
	}

	got := store.Get()
	expected := DefaultConfig()
	if got != expected {
		t.Errorf("expected defaults, got %+v", got)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestLoadCorruptedFile(t *testing.T) {
	path := tempConfigPath(t)
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
