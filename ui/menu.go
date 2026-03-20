package ui

import (
	"fyne.io/fyne/v2"

	"moment/core"
)

// ContextMenuDeps holds dependencies for the context menu.
type ContextMenuDeps struct {
	WindowMgr *core.WindowManager
	Config    *core.ConfigStore
	Quit      func()
}

// ContextMenu manages the system tray menu.
type ContextMenu struct {
	deps ContextMenuDeps
}

func NewContextMenu(deps ContextMenuDeps) *ContextMenu {
	return &ContextMenu{deps: deps}
}

// Menu returns a fyne.Menu (used for system tray).
func (c *ContextMenu) Menu() *fyne.Menu {
	return c.build()
}

func (c *ContextMenu) build() *fyne.Menu {
	return fyne.NewMenu("此刻",
		c.buildWindowMenu(),
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
