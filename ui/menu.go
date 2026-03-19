package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"moment/core"
)

// ContextMenuDeps holds the dependencies that the ContextMenu needs to
// interact with the core layer and other UI components.
type ContextMenuDeps struct {
	Clock     *ClockWidget
	Theme     *core.ThemeManager
	WindowMgr *core.WindowManager
	RestTimer *core.RestTimer
	Config    *core.ConfigStore
	Quit      func()
}

// ContextMenu builds and manages the right-click popup menu for the clock widget.
type ContextMenu struct {
	deps ContextMenuDeps
	menu *fyne.Menu
}

// NewContextMenu creates a ContextMenu wired to the given dependencies.
func NewContextMenu(deps ContextMenuDeps) *ContextMenu {
	cm := &ContextMenu{deps: deps}
	cm.menu = cm.build()
	return cm
}

// Menu returns the underlying fyne.Menu for use with widget.PopUpMenu.
func (c *ContextMenu) Menu() *fyne.Menu {
	return c.menu
}

// ShowAtPosition displays the context menu as a popup at the given position on the canvas.
func (c *ContextMenu) ShowAtPosition(cv fyne.Canvas, pos fyne.Position) {
	// Rebuild to reflect current state (checked items, etc.)
	c.menu = c.build()
	popup := widget.NewPopUpMenu(c.menu, cv)
	popup.ShowAtPosition(pos)
}

// build constructs the full menu tree.
func (c *ContextMenu) build() *fyne.Menu {
	return fyne.NewMenu("",
		c.displayModeItems()...,
	)
}

// displayModeItems returns all top-level menu items including sub-menus.
func (c *ContextMenu) displayModeItems() []*fyne.MenuItem {
	items := []*fyne.MenuItem{
		c.buildDisplayModeMenu(),
		c.buildThemeMenu(),
		c.buildWindowLevelMenu(),
		c.buildLockItem(),
		fyne.NewMenuItemSeparator(),
		c.buildRestMenu(),
		fyne.NewMenuItemSeparator(),
		c.buildQuitItem(),
	}
	return items
}

// buildDisplayModeMenu creates the "显示模式" sub-menu.
func (c *ContextMenu) buildDisplayModeMenu() *fyne.MenuItem {
	currentMode := core.ModeDigital
	if c.deps.Clock != nil {
		currentMode = c.deps.Clock.GetMode()
	}

	digital := fyne.NewMenuItem("数字时间", func() {
		if c.deps.Clock != nil {
			c.deps.Clock.SetMode(core.ModeDigital)
		}
		if c.deps.Config != nil {
			_ = c.deps.Config.Update(func(cfg *core.Config) {
				cfg.DisplayMode = core.ModeDigital
			})
		}
	})
	digital.Checked = currentMode == core.ModeDigital

	analog := fyne.NewMenuItem("模拟时钟", func() {
		if c.deps.Clock != nil {
			c.deps.Clock.SetMode(core.ModeAnalog)
		}
		if c.deps.Config != nil {
			_ = c.deps.Config.Update(func(cfg *core.Config) {
				cfg.DisplayMode = core.ModeAnalog
			})
		}
	})
	analog.Checked = currentMode == core.ModeAnalog

	timestamp := fyne.NewMenuItem("时间戳", func() {
		if c.deps.Clock != nil {
			c.deps.Clock.SetMode(core.ModeTimestamp)
		}
		if c.deps.Config != nil {
			_ = c.deps.Config.Update(func(cfg *core.Config) {
				cfg.DisplayMode = core.ModeTimestamp
			})
		}
	})
	timestamp.Checked = currentMode == core.ModeTimestamp

	item := fyne.NewMenuItem("显示模式", nil)
	item.ChildMenu = fyne.NewMenu("", digital, analog, timestamp)
	return item
}

// buildThemeMenu creates the "主题" sub-menu.
func (c *ContextMenu) buildThemeMenu() *fyne.MenuItem {
	currentMode := core.ThemeSystem
	if c.deps.Theme != nil {
		currentMode = c.deps.Theme.GetMode()
	}

	light := fyne.NewMenuItem("浅色", func() {
		if c.deps.Theme != nil {
			c.deps.Theme.SetMode(core.ThemeLight)
		}
		if c.deps.Config != nil {
			_ = c.deps.Config.Update(func(cfg *core.Config) {
				cfg.ThemeMode = core.ThemeLight
			})
		}
	})
	light.Checked = currentMode == core.ThemeLight

	dark := fyne.NewMenuItem("深色", func() {
		if c.deps.Theme != nil {
			c.deps.Theme.SetMode(core.ThemeDark)
		}
		if c.deps.Config != nil {
			_ = c.deps.Config.Update(func(cfg *core.Config) {
				cfg.ThemeMode = core.ThemeDark
			})
		}
	})
	dark.Checked = currentMode == core.ThemeDark

	system := fyne.NewMenuItem("跟随系统", func() {
		if c.deps.Theme != nil {
			c.deps.Theme.SetMode(core.ThemeSystem)
		}
		if c.deps.Config != nil {
			_ = c.deps.Config.Update(func(cfg *core.Config) {
				cfg.ThemeMode = core.ThemeSystem
			})
		}
	})
	system.Checked = currentMode == core.ThemeSystem

	item := fyne.NewMenuItem("主题", nil)
	item.ChildMenu = fyne.NewMenu("", light, dark, system)
	return item
}

// buildWindowLevelMenu creates the "窗口层级" sub-menu.
func (c *ContextMenu) buildWindowLevelMenu() *fyne.MenuItem {
	currentLevel := core.LevelTopMost
	if c.deps.WindowMgr != nil {
		currentLevel = c.deps.WindowMgr.GetLevel()
	}

	topMost := fyne.NewMenuItem("置顶", func() {
		if c.deps.WindowMgr != nil {
			c.deps.WindowMgr.SetLevel(core.LevelTopMost)
		}
	})
	topMost.Checked = currentLevel == core.LevelTopMost

	normal := fyne.NewMenuItem("普通", func() {
		if c.deps.WindowMgr != nil {
			c.deps.WindowMgr.SetLevel(core.LevelNormal)
		}
	})
	normal.Checked = currentLevel == core.LevelNormal

	item := fyne.NewMenuItem("窗口层级", nil)
	item.ChildMenu = fyne.NewMenu("", topMost, normal)
	return item
}

// buildLockItem creates the "锁定位置" / "解锁位置" toggle item.
func (c *ContextMenu) buildLockItem() *fyne.MenuItem {
	locked := false
	if c.deps.WindowMgr != nil {
		locked = c.deps.WindowMgr.IsLocked()
	}

	label := "锁定位置"
	if locked {
		label = "解锁位置"
	}

	item := fyne.NewMenuItem(label, func() {
		if c.deps.WindowMgr != nil {
			c.deps.WindowMgr.SetLocked(!locked)
		}
	})
	item.Checked = locked
	return item
}

// buildRestMenu creates the "休息提醒" sub-menu.
func (c *ContextMenu) buildRestMenu() *fyne.MenuItem {
	enabled := true
	if c.deps.RestTimer != nil {
		enabled = c.deps.RestTimer.IsEnabled()
	}

	toggleLabel := "关闭提醒"
	if !enabled {
		toggleLabel = "开启提醒"
	}
	toggle := fyne.NewMenuItem(toggleLabel, func() {
		if c.deps.RestTimer != nil {
			c.deps.RestTimer.SetEnabled(!enabled)
			if c.deps.Config != nil {
				_ = c.deps.Config.Update(func(cfg *core.Config) {
					cfg.RestEnabled = !enabled
				})
			}
		}
	})

	// Interval options
	intervals := []int{15, 20, 30, 45, 60, 90}
	currentInterval := 45
	if c.deps.RestTimer != nil {
		currentInterval = int(c.deps.RestTimer.GetInterval().Minutes())
	}

	var intervalItems []*fyne.MenuItem
	for _, mins := range intervals {
		m := mins // capture
		item := fyne.NewMenuItem(fmt.Sprintf("%d 分钟", m), func() {
			if c.deps.RestTimer != nil {
				c.deps.RestTimer.SetInterval(time.Duration(m) * time.Minute)
			}
			if c.deps.Config != nil {
				_ = c.deps.Config.Update(func(cfg *core.Config) {
					cfg.RestInterval = m
				})
			}
		})
		item.Checked = m == currentInterval
		intervalItems = append(intervalItems, item)
	}
	intervalMenu := fyne.NewMenuItem("提醒间隔", nil)
	intervalMenu.ChildMenu = fyne.NewMenu("", intervalItems...)

	// Opacity options
	opacities := []struct {
		label string
		value float64
	}{
		{"30%", 0.3},
		{"50%", 0.5},
		{"70%", 0.7},
		{"90%", 0.9},
	}
	currentOpacity := 0.7
	if c.deps.RestTimer != nil {
		currentOpacity = c.deps.RestTimer.GetMaxOpacity()
	}

	var opacityItems []*fyne.MenuItem
	for _, o := range opacities {
		op := o // capture
		item := fyne.NewMenuItem(op.label, func() {
			if c.deps.RestTimer != nil {
				c.deps.RestTimer.SetMaxOpacity(op.value)
			}
			if c.deps.Config != nil {
				_ = c.deps.Config.Update(func(cfg *core.Config) {
					cfg.RestOpacity = op.value
				})
			}
		})
		// Compare with small tolerance for float
		item.Checked = (op.value-currentOpacity) < 0.01 && (currentOpacity-op.value) < 0.01
		opacityItems = append(opacityItems, item)
	}
	opacityMenu := fyne.NewMenuItem("蒙版透明度", nil)
	opacityMenu.ChildMenu = fyne.NewMenu("", opacityItems...)

	item := fyne.NewMenuItem("休息提醒", nil)
	item.ChildMenu = fyne.NewMenu("", toggle, fyne.NewMenuItemSeparator(), intervalMenu, opacityMenu)
	return item
}

// buildQuitItem creates the "退出" menu item.
func (c *ContextMenu) buildQuitItem() *fyne.MenuItem {
	return fyne.NewMenuItem("退出", func() {
		if c.deps.Quit != nil {
			c.deps.Quit()
		}
	})
}
