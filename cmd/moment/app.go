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

// MomentApp is the top-level application controller.
type MomentApp struct {
	fyneApp   fyne.App
	window    fyne.Window
	clock     *ui.ClockWidget
	menu      *ui.ContextMenu
	windowMgr *core.WindowManager
	config    *core.ConfigStore
}

// NewMomentApp creates and fully initialises the application.
func NewMomentApp() *MomentApp {
	m := &MomentApp{}

	m.fyneApp = app.NewWithID("com.moment.clock")

	cs, err := core.NewConfigStore()
	if err != nil {
		log.Printf("config store init error: %v, using defaults", err)
	} else {
		m.config = cs
		if err := cs.Load(); err != nil {
			log.Printf("config load error: %v", err)
		}
	}

	// Main window — borderless
	if drv, ok := m.fyneApp.Driver().(interface {
		CreateSplashWindow() fyne.Window
	}); ok {
		m.window = drv.CreateSplashWindow()
	} else {
		m.window = m.fyneApp.NewWindow("此刻 Moment")
	}
	m.window.SetPadded(false)
	m.window.Resize(fyne.NewSize(300, 90))

	m.windowMgr = core.NewWindowManager(m.window, m.config)
	m.clock = ui.NewClockWidget()

	m.menu = ui.NewContextMenu(ui.ContextMenuDeps{
		WindowMgr: m.windowMgr,
		Config:    m.config,
		Quit:      m.Quit,
	})

	wrapper := newDraggableClock(m.clock, m.windowMgr)
	m.window.SetContent(wrapper)

	// System tray
	m.fyneApp.SetIcon(assets.IconResource)
	m.window.SetIcon(assets.IconResource)
	if desk, ok := m.fyneApp.(interface {
		SetSystemTrayMenu(menu *fyne.Menu)
		SetSystemTrayIcon(resource fyne.Resource)
	}); ok {
		desk.SetSystemTrayIcon(assets.IconResource)
		desk.SetSystemTrayMenu(m.menu.Menu())
	}

	m.window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			m.Quit()
		}
	})

	return m
}

func (m *MomentApp) Run() {
	go func() {
		time.Sleep(200 * time.Millisecond)
		core.RemoveTitleBar()
	}()
	m.window.ShowAndRun()
}

func (m *MomentApp) Quit() {
	if m.clock != nil {
		m.clock.Stop()
	}
	m.fyneApp.Quit()
}
