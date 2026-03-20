package ui

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// ClockWidget displays date + HH:MM:SS in flip-card style.
type ClockWidget struct {
	widget.BaseWidget

	mu          sync.RWMutex
	currentTime time.Time
	ticker      *time.Ticker
	done        chan struct{}
	onTick      func() // called every second so the parent can refresh
}

func NewClockWidget() *ClockWidget {
	c := &ClockWidget{
		currentTime: time.Now(),
		done:        make(chan struct{}),
	}
	c.ExtendBaseWidget(c)
	c.startTicker()
	return c
}

// SetOnTick registers a callback invoked every second after time updates.
// The parent widget should call its own Refresh() here.
func (c *ClockWidget) SetOnTick(fn func()) {
	c.mu.Lock()
	c.onTick = fn
	c.mu.Unlock()
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
	return newFlipCardRenderer(c)
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

// ===================== Colors =====================

var (
	bgDark   = color.NRGBA{R: 30, G: 30, B: 35, A: 255}
	cardDark = color.NRGBA{R: 55, G: 55, B: 62, A: 255}
	cardText = color.NRGBA{R: 240, G: 240, B: 245, A: 255}
	dateDim  = color.NRGBA{R: 160, G: 165, B: 175, A: 255}
	dotColor = color.NRGBA{R: 120, G: 130, B: 145, A: 255}
)

// ===================== Flip-Card Renderer =====================

const (
	flipCardW  float32 = 38
	flipCardH  float32 = 54
	flipGap    float32 = 3
	flipColonW float32 = 14
	flipPairG  float32 = 5
	flipFontSz float32 = 32
	flipCorner float32 = 7
	flipDotR   float32 = 3.5
	dateFontSz float32 = 12
	dateH      float32 = 18
	padY       float32 = 6
)

type digitCard struct {
	bg   *canvas.Rectangle
	text *canvas.Text
}

func newDigitCard() digitCard {
	bg := canvas.NewRectangle(cardDark)
	bg.CornerRadius = flipCorner
	txt := canvas.NewText("0", cardText)
	txt.TextSize = flipFontSz
	txt.Alignment = fyne.TextAlignCenter
	txt.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	return digitCard{bg: bg, text: txt}
}

type flipCardRenderer struct {
	clock     *ClockWidget
	bg        *canvas.Rectangle
	digits    [6]digitCard
	colonDots [4]*canvas.Circle
	dateLbl   *canvas.Text
	objects   []fyne.CanvasObject
}

func newFlipCardRenderer(clock *ClockWidget) *flipCardRenderer {
	r := &flipCardRenderer{clock: clock}
	r.bg = canvas.NewRectangle(bgDark)
	for i := 0; i < 6; i++ {
		r.digits[i] = newDigitCard()
	}
	for i := 0; i < 4; i++ {
		r.colonDots[i] = canvas.NewCircle(dotColor)
	}
	r.dateLbl = canvas.NewText("", dateDim)
	r.dateLbl.TextSize = dateFontSz
	r.dateLbl.Alignment = fyne.TextAlignCenter
	r.dateLbl.TextStyle = fyne.TextStyle{Monospace: true}
	r.buildObjects()
	r.Refresh()
	return r
}

func (r *flipCardRenderer) buildObjects() {
	r.objects = []fyne.CanvasObject{r.bg, r.dateLbl}
	for i := 0; i < 6; i++ {
		r.objects = append(r.objects, r.digits[i].bg, r.digits[i].text)
	}
	for i := 0; i < 4; i++ {
		r.objects = append(r.objects, r.colonDots[i])
	}
}

func (r *flipCardRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	r.dateLbl.Resize(fyne.NewSize(size.Width, dateH))
	r.dateLbl.Move(fyne.NewPos(0, padY-2))

	totalW := 6*flipCardW + 2*flipGap + 2*flipColonW + 4*flipPairG
	startX := (size.Width - totalW) / 2
	startY := dateH + padY

	x := startX
	for i := 0; i < 6; i++ {
		r.digits[i].bg.Resize(fyne.NewSize(flipCardW, flipCardH))
		r.digits[i].bg.Move(fyne.NewPos(x, startY))
		r.digits[i].text.Resize(fyne.NewSize(flipCardW, flipCardH))
		r.digits[i].text.Move(fyne.NewPos(x, startY))
		x += flipCardW
		if i == 1 || i == 3 {
			x += flipPairG
			ci := 0
			if i == 3 {
				ci = 2
			}
			dotCX := x + flipColonW/2
			r.colonDots[ci].Resize(fyne.NewSize(flipDotR*2, flipDotR*2))
			r.colonDots[ci].Move(fyne.NewPos(dotCX-flipDotR, startY+flipCardH*0.30-flipDotR))
			r.colonDots[ci+1].Resize(fyne.NewSize(flipDotR*2, flipDotR*2))
			r.colonDots[ci+1].Move(fyne.NewPos(dotCX-flipDotR, startY+flipCardH*0.70-flipDotR))
			x += flipColonW + flipPairG
		} else if i%2 == 0 {
			x += flipGap
		}
	}
}

func (r *flipCardRenderer) MinSize() fyne.Size {
	totalW := 6*flipCardW + 2*flipGap + 2*flipColonW + 4*flipPairG + 16
	return fyne.NewSize(totalW, flipCardH+dateH+padY*2+4)
}

func (r *flipCardRenderer) Refresh() {
	r.clock.mu.RLock()
	t := r.clock.currentTime
	r.clock.mu.RUnlock()

	r.dateLbl.Text = t.Format("2006-01-02 Mon")
	r.dateLbl.Color = dateDim
	canvas.Refresh(r.dateLbl)

	timeStr := t.Format("150405")
	for i := 0; i < 6 && i < len(timeStr); i++ {
		r.digits[i].text.Text = string(timeStr[i])
		r.digits[i].text.Color = cardText
		canvas.Refresh(r.digits[i].text)
	}
}

func (r *flipCardRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *flipCardRenderer) Destroy()                     {}
