package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"moment/assets"
	"moment/core"
	"moment/ui"
)

// MomentApp is the top-level application controller that initialises and
// wires together every component of the Moment desktop clock.
type MomentApp struct {
	fyneApp   fyne.App
	window    fyne.Window
	clock     *ui.ClockWidget
	menu      *ui.ContextMenu
	overlay   *ui.RestOverlay
	theme     *core.ThemeManager
	windowMgr *core.WindowManager
	restTimer *core.RestTimer
	config    *core.ConfigStore
}

// NewMomentApp creates and fully initialises the application.
func NewMomentApp() *MomentApp {
	m := &MomentApp{}

	// 1. Fyne app
	m.fyneApp = app.NewWithID("com.moment.clock")

	// 2. Config store
	cs, err := core.NewConfigStore()
	if err != nil {
		log.Printf("config store init error: %v, using defaults", err)
	} else {
		m.config = cs
		if err := cs.Load(); err != nil {
			log.Printf("config load error: %v", err)
		}
	}

	cfg := core.DefaultConfig()
	if m.config != nil {
		cfg = m.config.Get()
	}

	// 3. Theme manager
	m.theme = core.NewThemeManager(m.fyneApp)
	m.theme.SetMode(cfg.ThemeMode)

	// 4. Main window — borderless splash window
	m.window = m.fyneApp.NewWindow("此刻 Moment")
	m.window.SetPadded(false)
	m.window.Resize(fyne.NewSize(200, 80))

	// 5. Window manager (restores level, lock, position from config)
	m.windowMgr = core.NewWindowManager(m.window, m.config)

	// 6. Clock widget
	m.clock = ui.NewClockWidget(cfg.DisplayMode)

	// 7. Rest overlay + timer
	m.overlay = ui.NewRestOverlay(m.fyneApp)
	m.overlay.SetMaxOpacity(cfg.RestOpacity)
	m.overlay.SetOnDismiss(func() {
		if m.restTimer != nil {
			m.restTimer.Reset()
		}
	})

	m.restTimer = core.NewRestTimer(func() {
		opacity := m.restTimer.GetMaxOpacity()
		m.overlay.Show(opacity)
	})
	m.restTimer.SetMaxOpacity(cfg.RestOpacity)
	if cfg.RestInterval > 0 {
		m.restTimer.SetInterval(time.Duration(cfg.RestInterval) * time.Minute)
	}
	m.restTimer.SetEnabled(cfg.RestEnabled)

	// 8. Context menu
	m.menu = ui.NewContextMenu(ui.ContextMenuDeps{
		Clock:     m.clock,
		Theme:     m.theme,
		WindowMgr: m.windowMgr,
		RestTimer: m.restTimer,
		Config:    m.config,
		Quit:      m.Quit,
	})

	// Wrap clock in draggable container with right-click support
	wrapper := newDraggableClock(m.clock, m.windowMgr, m.window)
	wrapper.setMenu(m.menu)
	m.window.SetContent(wrapper)

	// 9. System tray — same menu as right-click, with the peach icon
	m.fyneApp.SetIcon(assets.IconResource)
	m.window.SetIcon(assets.IconResource)
	if desk, ok := m.fyneApp.(interface {
		SetSystemTrayMenu(menu *fyne.Menu)
		SetSystemTrayIcon(resource fyne.Resource)
	}); ok {
		desk.SetSystemTrayIcon(assets.IconResource)
		desk.SetSystemTrayMenu(m.menu.Menu())
	}

	// Escape key to quit
	m.window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			m.Quit()
		}
	})

	return m
}

// Run starts the Fyne event loop. Blocks until the app exits.
func (m *MomentApp) Run() {
	m.window.ShowAndRun()
}

// Quit cleans up resources and exits the application.
func (m *MomentApp) Quit() {
	if m.clock != nil {
		m.clock.Stop()
	}
	if m.restTimer != nil {
		m.restTimer.Stop()
	}
	if m.theme != nil {
		m.theme.Stop()
	}
	m.fyneApp.Quit()
}
