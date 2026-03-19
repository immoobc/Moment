package ui

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultFadeInDur  = 3 * time.Second
	defaultFadeOutDur = 2 * time.Second
	fadeStepInterval  = 50 * time.Millisecond // ~20 fps for smooth animation
)

// overlayColor is the green tint used for the rest overlay.
var overlayColor = color.NRGBA{R: 34, G: 139, B: 34, A: 0} // forest green, alpha set dynamically

// RestOverlay displays a full-screen semi-transparent green overlay to remind
// the user to take a break. It supports fade-in / fade-out animations and
// dismisses on mouse movement or key press.
type RestOverlay struct {
	mu         sync.Mutex
	app        fyne.App
	window     fyne.Window
	bg         *canvas.Rectangle
	maxOpacity float64
	fadeInDur  time.Duration
	fadeOutDur time.Duration
	active     bool
	onDismiss  func()

	// stopFade cancels an in-progress fade animation.
	stopFade chan struct{}
}

// NewRestOverlay creates a RestOverlay. The overlay window is created lazily
// on the first call to Show.
func NewRestOverlay(app fyne.App) *RestOverlay {
	return &RestOverlay{
		app:        app,
		maxOpacity: 0.7,
		fadeInDur:  defaultFadeInDur,
		fadeOutDur: defaultFadeOutDur,
	}
}

// SetOnDismiss registers a callback invoked when the overlay is dismissed.
func (o *RestOverlay) SetOnDismiss(fn func()) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.onDismiss = fn
}

// SetMaxOpacity sets the peak opacity for future Show calls.
// Requirements 6.6.
func (o *RestOverlay) SetMaxOpacity(opacity float64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}
	o.maxOpacity = opacity
}

// GetMaxOpacity returns the configured maximum opacity.
func (o *RestOverlay) GetMaxOpacity() float64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.maxOpacity
}

// IsActive returns whether the overlay is currently displayed.
func (o *RestOverlay) IsActive() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.active
}

// Show displays the overlay with a fade-in animation up to maxOpacity.
// Requirements 6.1: gradually increase opacity to full-screen green overlay.
func (o *RestOverlay) Show(maxOpacity float64) {
	o.mu.Lock()
	if o.active {
		o.mu.Unlock()
		return
	}
	o.active = true

	if maxOpacity <= 0 {
		maxOpacity = o.maxOpacity
	}

	o.ensureWindowLocked()
	o.stopFade = make(chan struct{})
	stop := o.stopFade
	o.mu.Unlock()

	o.window.Show()
	o.fadeIn(maxOpacity, stop)
}

// Dismiss hides the overlay with a fade-out animation.
// Requirements 6.3: gradually decrease opacity and disappear.
func (o *RestOverlay) Dismiss() {
	o.mu.Lock()
	if !o.active {
		o.mu.Unlock()
		return
	}

	// Cancel any running fade-in.
	if o.stopFade != nil {
		select {
		case <-o.stopFade:
		default:
			close(o.stopFade)
		}
	}
	o.stopFade = make(chan struct{})
	stop := o.stopFade
	cb := o.onDismiss
	o.mu.Unlock()

	o.fadeOut(stop, func() {
		o.mu.Lock()
		o.active = false
		if o.window != nil {
			o.window.Hide()
		}
		o.mu.Unlock()

		if cb != nil {
			cb()
		}
	})
}

// ensureWindowLocked creates the overlay window if it doesn't exist yet.
// Caller must hold the mutex.
func (o *RestOverlay) ensureWindowLocked() {
	if o.window != nil {
		return
	}

	w := o.app.NewWindow("Moment Rest Overlay")
	w.SetFullScreen(true)
	w.SetPadded(false)

	bg := canvas.NewRectangle(overlayColor)
	w.SetContent(bg)
	o.bg = bg
	o.window = w

	// Listen for mouse movement to dismiss — Fyne desktop.Hoverable on the
	// canvas object isn't ideal for full-screen detection, so we use a
	// custom canvas object that implements desktop.Hoverable.
	dismissObj := newDismissListener(o)
	w.SetContent(dismissObj)

	// Also listen for key presses.
	w.Canvas().SetOnTypedKey(func(_ *fyne.KeyEvent) {
		o.Dismiss()
	})
}

// fadeIn gradually increases the background alpha from 0 to target over fadeInDur.
func (o *RestOverlay) fadeIn(target float64, stop chan struct{}) {
	o.mu.Lock()
	dur := o.fadeInDur
	o.mu.Unlock()

	steps := int(dur / fadeStepInterval)
	if steps < 1 {
		steps = 1
	}

	for i := 1; i <= steps; i++ {
		select {
		case <-stop:
			return
		default:
		}
		alpha := target * float64(i) / float64(steps)
		o.setAlpha(alpha)
		time.Sleep(fadeStepInterval)
	}
	o.setAlpha(target)
}

// fadeOut gradually decreases the background alpha to 0 over fadeOutDur, then calls done.
func (o *RestOverlay) fadeOut(stop chan struct{}, done func()) {
	o.mu.Lock()
	dur := o.fadeOutDur
	o.mu.Unlock()

	currentAlpha := o.getAlpha()
	steps := int(dur / fadeStepInterval)
	if steps < 1 {
		steps = 1
	}

	for i := 1; i <= steps; i++ {
		select {
		case <-stop:
			if done != nil {
				done()
			}
			return
		default:
		}
		alpha := currentAlpha * (1.0 - float64(i)/float64(steps))
		o.setAlpha(alpha)
		time.Sleep(fadeStepInterval)
	}
	o.setAlpha(0)
	if done != nil {
		done()
	}
}

// setAlpha updates the overlay background color's alpha channel.
func (o *RestOverlay) setAlpha(a float64) {
	o.mu.Lock()
	bg := o.bg
	o.mu.Unlock()

	if bg == nil {
		return
	}
	bg.FillColor = color.NRGBA{
		R: overlayColor.R,
		G: overlayColor.G,
		B: overlayColor.B,
		A: uint8(a * 255),
	}
	bg.Refresh()
}

// getAlpha reads the current alpha from the background rectangle.
func (o *RestOverlay) getAlpha() float64 {
	o.mu.Lock()
	bg := o.bg
	o.mu.Unlock()

	if bg == nil {
		return 0
	}
	c, ok := bg.FillColor.(color.NRGBA)
	if !ok {
		return 0
	}
	return float64(c.A) / 255.0
}

// dismissListener is a full-size canvas object that detects mouse movement
// and triggers overlay dismissal. It also serves as the container for the
// background rectangle.
type dismissListener struct {
	widget.BaseWidget
	overlay *RestOverlay
	bg      *canvas.Rectangle
}

var _ desktop.Hoverable = (*dismissListener)(nil)

func newDismissListener(overlay *RestOverlay) *dismissListener {
	d := &dismissListener{
		overlay: overlay,
		bg:      overlay.bg,
	}
	d.ExtendBaseWidget(d)
	return d
}

// CreateRenderer returns a renderer that simply draws the background rect.
func (d *dismissListener) CreateRenderer() fyne.WidgetRenderer {
	return &dismissRenderer{bg: d.bg}
}

// MouseIn is called when the cursor enters the widget area.
func (d *dismissListener) MouseIn(_ *desktop.MouseEvent) {}

// MouseMoved is called when the cursor moves within the widget.
// Requirements 6.3: mouse movement triggers overlay dismissal.
func (d *dismissListener) MouseMoved(_ *desktop.MouseEvent) {
	d.overlay.Dismiss()
}

// MouseOut is called when the cursor leaves the widget area.
func (d *dismissListener) MouseOut() {}

// dismissRenderer renders the background rectangle at full size.
type dismissRenderer struct {
	bg *canvas.Rectangle
}

func (r *dismissRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))
}

func (r *dismissRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (r *dismissRenderer) Refresh() {
	r.bg.Refresh()
}

func (r *dismissRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg}
}

func (r *dismissRenderer) Destroy() {}
