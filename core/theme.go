package core

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// MomentTheme implements fyne.Theme with custom colors for Moment.
type MomentTheme struct {
	dark bool
}

var _ fyne.Theme = (*MomentTheme)(nil)

// Color returns the color for the given theme color name.
func (m *MomentTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Override variant based on our dark flag so the built-in theme
	// palette matches the user's choice.
	v := theme.VariantLight
	if m.dark {
		v = theme.VariantDark
	}
	return theme.DefaultTheme().Color(name, v)
}

// Font delegates to the default Fyne theme.
func (m *MomentTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon delegates to the default Fyne theme.
func (m *MomentTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size delegates to the default Fyne theme.
func (m *MomentTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// ThemeManager manages theme mode switching and system theme detection.
type ThemeManager struct {
	mu         sync.RWMutex
	mode       ThemeMode
	fyneApp    fyne.App
	onChange   func(ThemeMode)
	stopWatch  chan struct{}
	lightTheme *MomentTheme
	darkTheme  *MomentTheme
}

// NewThemeManager creates a ThemeManager with the given Fyne app.
// The initial mode defaults to ThemeSystem.
func NewThemeManager(app fyne.App) *ThemeManager {
	tm := &ThemeManager{
		mode:       ThemeSystem,
		fyneApp:    app,
		lightTheme: &MomentTheme{dark: false},
		darkTheme:  &MomentTheme{dark: true},
	}
	return tm
}

// SetMode changes the active theme mode and applies it.
func (t *ThemeManager) SetMode(mode ThemeMode) {
	t.mu.Lock()
	old := t.mode
	t.mode = mode

	// Stop any existing system watcher when switching away from ThemeSystem.
	if old == ThemeSystem && mode != ThemeSystem {
		t.stopSystemWatchLocked()
	}
	t.mu.Unlock()

	t.applyTheme()

	// Start system watcher if switching to ThemeSystem.
	if mode == ThemeSystem {
		t.startSystemWatch()
	}

	t.mu.RLock()
	cb := t.onChange
	t.mu.RUnlock()
	if cb != nil {
		cb(mode)
	}
}

// GetMode returns the current theme mode.
func (t *ThemeManager) GetMode() ThemeMode {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.mode
}

// CurrentTheme returns the fyne.Theme matching the current mode.
func (t *ThemeManager) CurrentTheme() fyne.Theme {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.resolveThemeLocked()
}

// SetOnChange registers a callback invoked after the mode changes.
func (t *ThemeManager) SetOnChange(fn func(ThemeMode)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onChange = fn
}

// Stop cleans up resources (stops the system theme watcher).
func (t *ThemeManager) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopSystemWatchLocked()
}

// resolveThemeLocked returns the concrete theme for the current mode.
// Caller must hold at least a read lock.
func (t *ThemeManager) resolveThemeLocked() fyne.Theme {
	switch t.mode {
	case ThemeLight:
		return t.lightTheme
	case ThemeDark:
		return t.darkTheme
	default: // ThemeSystem
		if t.isSystemDark() {
			return t.darkTheme
		}
		return t.lightTheme
	}
}

// applyTheme sets the resolved theme on the Fyne app.
func (t *ThemeManager) applyTheme() {
	if t.fyneApp == nil {
		return
	}
	th := t.CurrentTheme()
	t.fyneApp.Settings().SetTheme(th)
}

// isSystemDark attempts to detect whether the OS is using a dark theme.
// Fyne exposes the current system variant through its settings.
func (t *ThemeManager) isSystemDark() bool {
	if t.fyneApp == nil {
		return false
	}
	return t.fyneApp.Settings().ThemeVariant() == theme.VariantDark
}

// startSystemWatch begins polling the system theme so we can react to changes.
func (t *ThemeManager) startSystemWatch() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Already watching.
	if t.stopWatch != nil {
		return
	}

	stop := make(chan struct{})
	t.stopWatch = stop

	go func() {
		listener := make(chan fyne.Settings)
		t.fyneApp.Settings().AddChangeListener(listener)
		for {
			select {
			case <-stop:
				return
			case <-listener:
				t.mu.RLock()
				mode := t.mode
				t.mu.RUnlock()
				if mode == ThemeSystem {
					t.applyTheme()
				}
			}
		}
	}()
}

// stopSystemWatchLocked stops the system theme watcher goroutine.
// Caller must hold the write lock.
func (t *ThemeManager) stopSystemWatchLocked() {
	if t.stopWatch != nil {
		close(t.stopWatch)
		t.stopWatch = nil
	}
}
