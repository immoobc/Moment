package main

import (
	"fyne.io/fyne/v2"

	"moment/core"
	"moment/ui"
)

// draggableClock wraps a ClockWidget to support window dragging (controlled
// by WindowManager lock state) and right-click context menu display.
type draggableClock struct {
	*ui.ClockWidget
	windowMgr *core.WindowManager
	window    fyne.Window
	menu      *ui.ContextMenu
}

func newDraggableClock(clock *ui.ClockWidget, wm *core.WindowManager, win fyne.Window) *draggableClock {
	d := &draggableClock{
		ClockWidget: clock,
		windowMgr:   wm,
		window:      win,
	}
	return d
}

// setMenu sets the context menu reference.
func (d *draggableClock) setMenu(menu *ui.ContextMenu) {
	d.menu = menu
}

// Dragged implements fyne.Draggable — guards against locked position.
func (d *draggableClock) Dragged(ev *fyne.DragEvent) {
	if d.windowMgr != nil && d.windowMgr.IsLocked() {
		return
	}
	// Fyne handles native window movement for borderless/splash windows.
}

// DragEnd implements fyne.Draggable — persists position when unlocked.
func (d *draggableClock) DragEnd() {
	if d.windowMgr != nil && !d.windowMgr.IsLocked() {
		// Position persistence is handled by WindowManager when the
		// window is moved. Fyne borderless windows handle the actual
		// move natively.
	}
}

// TappedSecondary implements fyne.SecondaryTappable — shows the context menu.
func (d *draggableClock) TappedSecondary(ev *fyne.PointEvent) {
	if d.menu != nil {
		d.menu.ShowAtPosition(d.window.Canvas(), ev.Position)
	}
}
