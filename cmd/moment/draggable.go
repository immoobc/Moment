package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"moment/core"
	"moment/ui"
)

// draggableClock wraps a ClockWidget to support window dragging.
type draggableClock struct {
	widget.BaseWidget
	clock      *ui.ClockWidget
	windowMgr  *core.WindowManager
	dragActive bool
}

func newDraggableClock(clock *ui.ClockWidget, wm *core.WindowManager) *draggableClock {
	d := &draggableClock{
		clock:     clock,
		windowMgr: wm,
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
