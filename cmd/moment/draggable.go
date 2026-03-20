package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"moment/core"
	"moment/ui"
)

// draggableClock wraps a ClockWidget to support window dragging and right-click menu.
type draggableClock struct {
	widget.BaseWidget
	clock      *ui.ClockWidget
	windowMgr  *core.WindowManager
	dragActive bool
	menu       *ui.ContextMenu
	window     fyne.Window
}

func newDraggableClock(clock *ui.ClockWidget, wm *core.WindowManager, menu *ui.ContextMenu, win fyne.Window) *draggableClock {
	d := &draggableClock{
		clock:     clock,
		windowMgr: wm,
		menu:      menu,
		window:    win,
	}
	d.ExtendBaseWidget(d)
	clock.SetOnTick(func() {
		d.Refresh()
	})
	return d
}

func (d *draggableClock) CreateRenderer() fyne.WidgetRenderer {
	return d.clock.CreateRenderer()
}

func (d *draggableClock) Dragged(_ *fyne.DragEvent) {
	if d.windowMgr == nil || d.windowMgr.IsLocked() {
		return
	}
	if !d.dragActive {
		d.dragActive = true
		d.windowMgr.BeginDrag()
	}
	d.windowMgr.DragUpdate()
}

func (d *draggableClock) DragEnd() {
	d.dragActive = false
	if d.windowMgr != nil && !d.windowMgr.IsLocked() {
		d.windowMgr.DragEnd()
	}
}

// TappedSecondary handles right-click to show popup menu.
func (d *draggableClock) TappedSecondary(ev *fyne.PointEvent) {
	if d.menu == nil || d.window == nil {
		return
	}
	m := d.menu.PopupMenu()
	popup := widget.NewPopUpMenu(m, d.window.Canvas())
	popup.ShowAtPosition(ev.AbsolutePosition)
}
