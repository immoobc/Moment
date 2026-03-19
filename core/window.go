package core

import (
	"sync"

	"fyne.io/fyne/v2"
)

// WindowManager manages window z-order level, position locking, and position persistence.
type WindowManager struct {
	mu       sync.RWMutex
	window   fyne.Window
	level    WindowLevel
	locked   bool
	position fyne.Position
	config   *ConfigStore
}

// NewWindowManager creates a WindowManager for the given window.
// If a ConfigStore is provided, position changes are persisted automatically.
func NewWindowManager(window fyne.Window, config *ConfigStore) *WindowManager {
	wm := &WindowManager{
		window: window,
		level:  LevelTopMost,
		locked: false,
		config: config,
	}
	if config != nil {
		cfg := config.Get()
		wm.level = cfg.WindowLevel
		wm.locked = cfg.Locked
		wm.position = fyne.NewPos(cfg.PositionX, cfg.PositionY)
	}
	return wm
}

// SetLevel changes the window z-order level and applies it immediately.
func (w *WindowManager) SetLevel(level WindowLevel) {
	w.mu.Lock()
	w.level = level
	w.mu.Unlock()

	w.applyLevel()

	if w.config != nil {
		_ = w.config.Update(func(c *Config) {
			c.WindowLevel = level
		})
	}
}

// GetLevel returns the current window level.
func (w *WindowManager) GetLevel() WindowLevel {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.level
}

// SetLocked sets the position lock state.
// When locked, drag operations should be rejected by the UI layer.
func (w *WindowManager) SetLocked(locked bool) {
	w.mu.Lock()
	w.locked = locked
	w.mu.Unlock()

	if w.config != nil {
		_ = w.config.Update(func(c *Config) {
			c.Locked = locked
		})
	}
}

// IsLocked returns whether the window position is locked.
func (w *WindowManager) IsLocked() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.locked
}

// SetPosition updates the stored window position.
// The position is only accepted when the window is unlocked.
// Returns true if the position was accepted.
func (w *WindowManager) SetPosition(pos fyne.Position) bool {
	w.mu.Lock()
	if w.locked {
		w.mu.Unlock()
		return false
	}
	w.position = pos
	w.mu.Unlock()

	if w.config != nil {
		_ = w.config.Update(func(c *Config) {
			c.PositionX = pos.X
			c.PositionY = pos.Y
		})
	}
	return true
}

// GetPosition returns the current stored window position.
func (w *WindowManager) GetPosition() fyne.Position {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.position
}

// applyLevel applies the current window level using platform-specific APIs.
func (w *WindowManager) applyLevel() {
	w.mu.RLock()
	level := w.level
	win := w.window
	w.mu.RUnlock()

	if win == nil {
		return
	}

	applyWindowLevel(win, level)
}
