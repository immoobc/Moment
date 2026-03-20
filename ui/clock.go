package ui

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"moment/core"
)

// ClockWidget displays date + HH:MM:SS with theme support.
type ClockWidget struct {
	widget.BaseWidget

	mu          sync.RWMutex
	currentTime time.Time
	theme       core.ThemeMode
	ticker      *time.Ticker
	done        chan struct{}
	onTick      func()
}

func NewClockWidget(theme core.ThemeMode) *ClockWidget {
	c := &ClockWidget{
		currentTime: time.Now(),
		theme:       theme,
		done:        make(chan struct{}),
	}
	c.ExtendBaseWidget(c)
	c.startTicker()
	return c
}

func (c *ClockWidget) SetOnTick(fn func()) {
	c.mu.Lock()
	c.onTick = fn
	c.mu.Unlock()
}

// SetTheme switches the visual theme and triggers a full refresh.
func (c *ClockWidget) SetTheme(t core.ThemeMode) {
	c.mu.Lock()
	c.theme = t
	c.mu.Unlock()
	c.Refresh()
}

func (c *ClockWidget) GetTheme() core.ThemeMode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.theme
}

func (c *ClockWidget) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	select {
	case <-c.done:
	default:
		close(c.done)
	}
}

func (c *ClockWidget) CreateRenderer() fyne.WidgetRenderer {
	return newCalendarRenderer(c)
}

func (c *ClockWidget) startTicker() {
	c.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-c.ticker.C:
				c.mu.Lock()
				c.currentTime = time.Now()
				cb := c.onTick
				c.mu.Unlock()
				c.Refresh()
				if cb != nil {
					cb()
				}
			}
		}
	}()
}

// ===================== Theme Color Palettes =====================

type colorPalette struct {
	bg       color.NRGBA
	header   color.NRGBA
	dateText color.NRGBA
	cardBg   color.NRGBA
	digitTxt color.NRGBA
	dotColor color.NRGBA
	weekday  color.NRGBA
}

var lightPalette = colorPalette{
	bg:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	header:   color.NRGBA{R: 220, G: 50, B: 47, A: 255},
	dateText: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	cardBg:   color.NRGBA{R: 245, G: 245, B: 247, A: 255},
	digitTxt: color.NRGBA{R: 50, G: 50, B: 55, A: 255},
	dotColor: color.NRGBA{R: 200, G: 60, B: 60, A: 255},
	weekday:  color.NRGBA{R: 130, G: 130, B: 135, A: 255},
}

var darkPalette = colorPalette{
	bg:       color.NRGBA{R: 30, G: 30, B: 35, A: 255},
	header:   color.NRGBA{R: 180, G: 40, B: 40, A: 255},
	dateText: color.NRGBA{R: 230, G: 230, B: 235, A: 255},
	cardBg:   color.NRGBA{R: 50, G: 50, B: 58, A: 255},
	digitTxt: color.NRGBA{R: 230, G: 230, B: 240, A: 255},
	dotColor: color.NRGBA{R: 200, G: 70, B: 70, A: 255},
	weekday:  color.NRGBA{R: 160, G: 160, B: 170, A: 255},
}

func paletteFor(t core.ThemeMode) colorPalette {
	if t == core.ThemeDark {
		return darkPalette
	}
	return lightPalette
}

// ===================== Layout Constants =====================

const (
	cardW      float32 = 34
	cardH      float32 = 46
	cardGap    float32 = 3
	colonW     float32 = 12
	pairGap    float32 = 4
	fontSize   float32 = 28
	cornerR    float32 = 6
	dotR       float32 = 3
	headerH    float32 = 32
	weekdayH   float32 = 18
	padTop     float32 = 6
	padBottom  float32 = 18
	padSide    float32 = 20
	headerCorR float32 = 12
)

// ===================== Calendar-Style Renderer =====================

type digitCard struct {
	bg   *canvas.Rectangle
	text *canvas.Text
}

func newDigitCard(p colorPalette) digitCard {
	bg := canvas.NewRectangle(p.cardBg)
	bg.CornerRadius = cornerR
	txt := canvas.NewText("0", p.digitTxt)
	txt.TextSize = fontSize
	txt.Alignment = fyne.TextAlignCenter
	txt.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	return digitCard{bg: bg, text: txt}
}

type calendarRenderer struct {
	clock   *ClockWidget
	bg      *canvas.Rectangle
	header  *canvas.Rectangle
	dateLbl *canvas.Text
	weekLbl *canvas.Text
	digits  [6]digitCard
	dots    [4]*canvas.Circle
	objects []fyne.CanvasObject
}

func newCalendarRenderer(clock *ClockWidget) *calendarRenderer {
	p := paletteFor(clock.GetTheme())
	r := &calendarRenderer{clock: clock}

	r.bg = canvas.NewRectangle(p.bg)
	r.bg.CornerRadius = headerCorR

	r.header = canvas.NewRectangle(p.header)
	r.header.CornerRadius = headerCorR

	r.dateLbl = canvas.NewText("", p.dateText)
	r.dateLbl.TextSize = 12
	r.dateLbl.Alignment = fyne.TextAlignCenter
	r.dateLbl.TextStyle = fyne.TextStyle{Bold: true}

	r.weekLbl = canvas.NewText("", p.weekday)
	r.weekLbl.TextSize = 11
	r.weekLbl.Alignment = fyne.TextAlignCenter

	for i := range 6 {
		r.digits[i] = newDigitCard(p)
	}
	for i := range 4 {
		r.dots[i] = canvas.NewCircle(p.dotColor)
	}

	r.buildObjects()
	r.Refresh()
	return r
}

func (r *calendarRenderer) buildObjects() {
	r.objects = []fyne.CanvasObject{r.bg, r.header, r.dateLbl, r.weekLbl}
	for i := range 6 {
		r.objects = append(r.objects, r.digits[i].bg, r.digits[i].text)
	}
	for i := range 4 {
		r.objects = append(r.objects, r.dots[i])
	}
}

func (r *calendarRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	r.header.Resize(fyne.NewSize(size.Width, headerH))
	r.header.Move(fyne.NewPos(0, 0))

	r.dateLbl.Resize(fyne.NewSize(size.Width, headerH))
	r.dateLbl.Move(fyne.NewPos(0, 8))

	r.weekLbl.Resize(fyne.NewSize(size.Width, weekdayH))
	r.weekLbl.Move(fyne.NewPos(0, headerH+2))

	digitsW := 6*cardW + 2*cardGap + 2*colonW + 4*pairGap
	startX := (size.Width - digitsW) / 2
	startY := headerH + weekdayH + padTop

	x := startX
	for i := range 6 {
		r.digits[i].bg.Resize(fyne.NewSize(cardW, cardH))
		r.digits[i].bg.Move(fyne.NewPos(x, startY))
		r.digits[i].text.Resize(fyne.NewSize(cardW, cardH))
		r.digits[i].text.Move(fyne.NewPos(x, startY))
		x += cardW
		if i == 1 || i == 3 {
			x += pairGap
			ci := 0
			if i == 3 {
				ci = 2
			}
			dotCX := x + colonW/2
			r.dots[ci].Resize(fyne.NewSize(dotR*2, dotR*2))
			r.dots[ci].Move(fyne.NewPos(dotCX-dotR, startY+cardH*0.30-dotR))
			r.dots[ci+1].Resize(fyne.NewSize(dotR*2, dotR*2))
			r.dots[ci+1].Move(fyne.NewPos(dotCX-dotR, startY+cardH*0.70-dotR))
			x += colonW + pairGap
		} else if i%2 == 0 {
			x += cardGap
		}
	}
}

func (r *calendarRenderer) MinSize() fyne.Size {
	totalW := 6*cardW + 2*cardGap + 2*colonW + 4*pairGap + 2*padSide
	totalH := headerH + weekdayH + padTop + cardH + padBottom
	return fyne.NewSize(totalW, totalH)
}

var weekdayNames = [7]string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

func (r *calendarRenderer) Refresh() {
	r.clock.mu.RLock()
	t := r.clock.currentTime
	theme := r.clock.theme
	r.clock.mu.RUnlock()

	p := paletteFor(theme)

	// Update colors for theme switch
	r.bg.FillColor = p.bg
	canvas.Refresh(r.bg)
	r.header.FillColor = p.header
	canvas.Refresh(r.header)

	r.dateLbl.Text = fmt.Sprintf("%d年%02d月%02d日", t.Year(), t.Month(), t.Day())
	r.dateLbl.Color = p.dateText
	canvas.Refresh(r.dateLbl)

	r.weekLbl.Text = weekdayNames[t.Weekday()]
	r.weekLbl.Color = p.weekday
	canvas.Refresh(r.weekLbl)

	timeStr := t.Format("150405")
	for i := range 6 {
		r.digits[i].bg.FillColor = p.cardBg
		canvas.Refresh(r.digits[i].bg)
		if i < len(timeStr) {
			r.digits[i].text.Text = string(timeStr[i])
		}
		r.digits[i].text.Color = p.digitTxt
		canvas.Refresh(r.digits[i].text)
	}
	for i := range 4 {
		r.dots[i].FillColor = p.dotColor
		canvas.Refresh(r.dots[i])
	}
}

func (r *calendarRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *calendarRenderer) Destroy()                     {}
