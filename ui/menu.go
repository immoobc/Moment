package ui

import (
	"fyne.io/fyne/v2"

	"moment/core"
)

// ContextMenuDeps holds dependencies for the context menu.
type ContextMenuDeps struct {
	WindowMgr   *core.WindowManager
	Config      *core.ConfigStore
	Quit        func()
	SwitchTheme func(core.ThemeMode) // callback to switch theme
}

// ContextMenu manages both the system tray menu and the right-click popup.
type ContextMenu struct {
	deps ContextMenuDeps
}

func NewContextMenu(deps ContextMenuDeps) *ContextMenu {
	return &ContextMenu{deps: deps}
}

// Menu returns a fyne.Menu for the system tray.
func (c *ContextMenu) Menu() *fyne.Menu {
	return c.build()
}

// PopupMenu returns a fyne.Menu suitable for right-click popup on the clock.
// It has the same items as the tray menu.
func (c *ContextMenu) PopupMenu() *fyne.Menu {
	return c.build()
}

func (c *ContextMenu) build() *fyne.Menu {
	return fyne.NewMenu("此刻",
		c.buildWindowMenu(),
		c.buildThemeMenu(),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("退出", func() {
			if c.deps.Quit != nil {
				c.deps.Quit()
			}
		}),
	)
}

func (c *ContextMenu) buildWindowMenu() *fyne.MenuItem {
	curLevel := core.LevelTopMost
	locked := false
	if c.deps.WindowMgr != nil {
		curLevel = c.deps.WindowMgr.GetLevel()
		locked = c.deps.WindowMgr.IsLocked()
	}

	top := fyne.NewMenuItem("置顶", func() {
		if c.deps.WindowMgr != nil {
			c.deps.WindowMgr.SetLevel(core.LevelTopMost)
		}
	})
	top.Checked = curLevel == core.LevelTopMost

	normal := fyne.NewMenuItem("普通层级", func() {
		if c.deps.WindowMgr != nil {
			c.deps.WindowMgr.SetLevel(core.LevelNormal)
		}
	})
	normal.Checked = curLevel == core.LevelNormal

	lockLabel := "🔒 锁定位置"
	if locked {
		lockLabel = "🔓 解锁位置"
	}
	lock := fyne.NewMenuItem(lockLabel, func() {
		if c.deps.WindowMgr != nil {
			c.deps.WindowMgr.SetLocked(!locked)
		}
	})

	item := fyne.NewMenuItem("窗口", nil)
	item.ChildMenu = fyne.NewMenu("", top, normal, fyne.NewMenuItemSeparator(), lock)
	return item
}

func (c *ContextMenu) buildThemeMenu() *fyne.MenuItem {
	curTheme := core.ThemeLight
	if c.deps.Config != nil {
		curTheme = c.deps.Config.Get().Theme
	}

	light := fyne.NewMenuItem("☀ 日历白", func() {
		if c.deps.SwitchTheme != nil {
			c.deps.SwitchTheme(core.ThemeLight)
		}
	})
	light.Checked = curTheme == core.ThemeLight

	dark := fyne.NewMenuItem("🌙 暗夜黑", func() {
		if c.deps.SwitchTheme != nil {
			c.deps.SwitchTheme(core.ThemeDark)
		}
	})
	dark.Checked = curTheme == core.ThemeDark

	item := fyne.NewMenuItem("主题", nil)
	item.ChildMenu = fyne.NewMenu("", light, dark)
	return item
}
