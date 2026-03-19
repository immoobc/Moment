package ui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"moment/core"
)

// ClockWidget is a Fyne widget that displays the current time in one of
// three modes: digital, analog, or Unix timestamp.
type ClockWidget struct {
	widget.BaseWidget

	mu          sync.RWMutex
	mode        core.DisplayMode
	currentTime time.Time
	ticker      *time.Ticker
	done        chan struct{}
}

// NewClockWidget creates a ClockWidget with the given initial display mode
// and starts a background ticker that refreshes every second.
func NewClockWidget(mode core.DisplayMode) *ClockWidget {
	c := &ClockWidget{
		mode:        mode,
		currentTime: time.Now(),
		done:        make(chan struct{}),
	}
	c.ExtendBaseWidget(c)
	c.startTicker()
	return c
}

// SetMode changes the display mode and refreshes the widget.
func (c *ClockWidget) SetMode(mode core.DisplayMode) {
	c.mu.Lock()
	c.mode = mode
	c.mu.Unlock()
	c.Refresh()
}

// GetMode returns the current display mode.
func (c *ClockWidget) GetMode() core.DisplayMode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode
}

// Stop stops the background ticker. Call this when the widget is no longer needed.
func (c *ClockWidget) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	select {
	case <-c.done:
		// already closed
	default:
		close(c.done)
	}
}

// CreateRenderer returns the appropriate renderer based on the current mode.
func (c *ClockWidget) CreateRenderer() fyne.WidgetRenderer {
	c.mu.RLock()
	mode := c.mode
	c.mu.RUnlock()

	if mode == core.ModeAnalog {
		return newAnalogClockRenderer(c)
	}
	return newTextClockRenderer(c)
}

// FormatTime returns the formatted time string for the given time and mode.
func FormatTime(t time.Time, mode core.DisplayMode) string {
	switch mode {
	case core.ModeDigital:
		return t.Format("15:04:05")
	case core.ModeTimestamp:
		return fmt.Sprintf("%d", t.Unix())
	default:
		return t.Format("15:04:05")
	}
}

// startTicker launches a goroutine that updates currentTime every second.
func (c *ClockWidget) startTicker() {
	c.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-c.done:
				return
			case t := <-c.ticker.C:
				c.mu.Lock()
				c.currentTime = t
				c.mu.Unlock()
				c.Refresh()
			}
		}
	}()
}

// textClockRenderer renders digital time and timestamp modes using a text label.
type textClockRenderer struct {
	clock *ClockWidget
	label *canvas.Text
}

func newTextClockRenderer(clock *ClockWidget) *textClockRenderer {
	clock.mu.RLock()
	text := FormatTime(clock.currentTime, clock.mode)
	clock.mu.RUnlock()

	label := canvas.NewText(text, nil)
	label.TextSize = 28
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Monospace: true}

	return &textClockRenderer{
		clock: clock,
		label: label,
	}
}

func (r *textClockRenderer) Layout(size fyne.Size) {
	r.label.Resize(size)
	r.label.Move(fyne.NewPos(0, 0))
}

func (r *textClockRenderer) MinSize() fyne.Size {
	return r.label.MinSize()
}

func (r *textClockRenderer) Refresh() {
	r.clock.mu.RLock()
	mode := r.clock.mode
	t := r.clock.currentTime
	r.clock.mu.RUnlock()

	// If mode switched to analog, the widget will recreate the renderer.
	if mode == core.ModeAnalog {
		return
	}

	r.label.Text = FormatTime(t, mode)
	r.label.Refresh()
}

func (r *textClockRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.label}
}

func (r *textClockRenderer) Destroy() {}
