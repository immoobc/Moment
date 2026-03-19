package ui

import (
	"image/color"
	"math"
	"moment/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const (
	analogMinSize = 120 // minimum width/height in dp
	markerCount   = 12
	hourHandRatio = 0.50 // fraction of radius
	minHandRatio  = 0.70
	secHandRatio  = 0.85
	markerInner   = 0.88 // inner end of tick mark
	markerOuter   = 0.98 // outer end of tick mark
)

// analogClockRenderer draws a classic round clock face with hour, minute and
// second hands using Fyne canvas primitives.
type analogClockRenderer struct {
	clock   *ClockWidget
	face    *canvas.Circle
	border  *canvas.Circle
	hour    *canvas.Line
	minute  *canvas.Line
	second  *canvas.Line
	center  *canvas.Circle
	markers [markerCount]*canvas.Line
	objects []fyne.CanvasObject
}

func newAnalogClockRenderer(clock *ClockWidget) *analogClockRenderer {
	face := canvas.NewCircle(color.Transparent)
	border := canvas.NewCircle(color.Transparent)
	border.StrokeColor = theme.Color(theme.ColorNameForeground)
	border.StrokeWidth = 2

	hourLine := canvas.NewLine(theme.Color(theme.ColorNameForeground))
	hourLine.StrokeWidth = 4

	minLine := canvas.NewLine(theme.Color(theme.ColorNameForeground))
	minLine.StrokeWidth = 2.5

	secLine := canvas.NewLine(color.NRGBA{R: 220, G: 50, B: 50, A: 255})
	secLine.StrokeWidth = 1.5

	centerDot := canvas.NewCircle(theme.Color(theme.ColorNameForeground))

	var markers [markerCount]*canvas.Line
	for i := 0; i < markerCount; i++ {
		m := canvas.NewLine(theme.Color(theme.ColorNameForeground))
		m.StrokeWidth = 1.5
		markers[i] = m
	}

	r := &analogClockRenderer{
		clock:   clock,
		face:    face,
		border:  border,
		hour:    hourLine,
		minute:  minLine,
		second:  secLine,
		center:  centerDot,
		markers: markers,
	}
	r.buildObjects()
	return r
}

func (r *analogClockRenderer) buildObjects() {
	r.objects = []fyne.CanvasObject{r.face, r.border}
	for _, m := range r.markers {
		r.objects = append(r.objects, m)
	}
	r.objects = append(r.objects, r.hour, r.minute, r.second, r.center)
}

func (r *analogClockRenderer) Layout(size fyne.Size) {
	side := size.Width
	if size.Height < side {
		side = size.Height
	}
	cx := size.Width / 2
	cy := size.Height / 2
	radius := side / 2

	// Face and border circles.
	r.face.Resize(fyne.NewSize(side, side))
	r.face.Move(fyne.NewPos(cx-radius, cy-radius))
	r.border.Resize(fyne.NewSize(side, side))
	r.border.Move(fyne.NewPos(cx-radius, cy-radius))

	// Center dot.
	dotR := float32(3)
	r.center.Resize(fyne.NewSize(dotR*2, dotR*2))
	r.center.Move(fyne.NewPos(cx-dotR, cy-dotR))

	// Tick markers.
	for i := 0; i < markerCount; i++ {
		angle := float64(i)*(2*math.Pi/markerCount) - math.Pi/2
		inner := float32(markerInner) * radius
		outer := float32(markerOuter) * radius
		sin, cos := math.Sincos(angle)
		r.markers[i].Position1 = fyne.NewPos(cx+inner*float32(cos), cy+inner*float32(sin))
		r.markers[i].Position2 = fyne.NewPos(cx+outer*float32(cos), cy+outer*float32(sin))
	}

	r.layoutHands(cx, cy, radius)
}

func (r *analogClockRenderer) layoutHands(cx, cy, radius float32) {
	r.clock.mu.RLock()
	t := r.clock.currentTime
	r.clock.mu.RUnlock()

	h, m, s := t.Hour()%12, t.Minute(), t.Second()

	// Angles (0 = 12 o'clock, clockwise).
	hourAngle := (float64(h)+float64(m)/60.0)*(2*math.Pi/12) - math.Pi/2
	minAngle := (float64(m)+float64(s)/60.0)*(2*math.Pi/60) - math.Pi/2
	secAngle := float64(s)*(2*math.Pi/60) - math.Pi/2

	setHand(r.hour, cx, cy, radius*float32(hourHandRatio), hourAngle)
	setHand(r.minute, cx, cy, radius*float32(minHandRatio), minAngle)
	setHand(r.second, cx, cy, radius*float32(secHandRatio), secAngle)
}

// HandAngles computes the hour, minute and second hand angles (in radians,
// 0 = 3 o'clock, counter-clockwise positive) for the given time.
// Exported for testing.
func HandAngles(h, m, s int) (hourAngle, minAngle, secAngle float64) {
	h = h % 12
	hourAngle = (float64(h)+float64(m)/60.0)*(2*math.Pi/12) - math.Pi/2
	minAngle = (float64(m)+float64(s)/60.0)*(2*math.Pi/60) - math.Pi/2
	secAngle = float64(s)*(2*math.Pi/60) - math.Pi/2
	return
}

func setHand(line *canvas.Line, cx, cy, length float32, angle float64) {
	sin, cos := math.Sincos(angle)
	line.Position1 = fyne.NewPos(cx, cy)
	line.Position2 = fyne.NewPos(cx+length*float32(cos), cy+length*float32(sin))
}

func (r *analogClockRenderer) MinSize() fyne.Size {
	return fyne.NewSize(analogMinSize, analogMinSize)
}

func (r *analogClockRenderer) Refresh() {
	r.clock.mu.RLock()
	mode := r.clock.mode
	r.clock.mu.RUnlock()

	// If mode switched away from analog, nothing to do here.
	if mode != core.ModeAnalog {
		return
	}

	// Update hand positions using the current widget size.
	size := r.clock.Size()
	side := size.Width
	if size.Height < side {
		side = size.Height
	}
	cx := size.Width / 2
	cy := size.Height / 2
	radius := side / 2

	r.layoutHands(cx, cy, radius)

	// Refresh colors in case theme changed.
	fg := theme.Color(theme.ColorNameForeground)
	r.border.StrokeColor = fg
	r.hour.StrokeColor = fg
	r.minute.StrokeColor = fg
	r.center.FillColor = fg

	r.border.Refresh()
	r.hour.Refresh()
	r.minute.Refresh()
	r.second.Refresh()
	r.center.Refresh()
}

func (r *analogClockRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *analogClockRenderer) Destroy() {}
