package core

import (
	"sync"
	"time"
)

// RestTimer manages the periodic rest-reminder cycle.
// When the configured interval elapses it invokes the onTrigger callback
// (typically used to show the RestOverlay). After the overlay is dismissed
// the caller should call Reset to restart the cycle.
type RestTimer struct {
	mu          sync.RWMutex
	interval    time.Duration // time between reminders (default 45 min)
	displayTime time.Duration // how long the overlay stays visible (default 5 min)
	maxOpacity  float64       // maximum overlay opacity 0.0–1.0
	enabled     bool
	running     bool
	timer       *time.Timer
	onTrigger   func()

	// newTimer is an indirection so tests can supply a fake.
	// Production code leaves it nil and the default time.NewTimer is used.
	newTimer func(d time.Duration) *time.Timer
}

// NewRestTimer creates a RestTimer with sensible defaults.
// onTrigger is called when the interval elapses.
func NewRestTimer(onTrigger func()) *RestTimer {
	return &RestTimer{
		interval:    45 * time.Minute,
		displayTime: 5 * time.Minute,
		maxOpacity:  0.7,
		enabled:     true,
		onTrigger:   onTrigger,
	}
}

// SetInterval updates the rest interval and restarts the timer if running.
// Requirements 6.5: new interval applied immediately.
func (r *RestTimer) SetInterval(d time.Duration) {
	r.mu.Lock()
	r.interval = d
	wasRunning := r.running
	r.mu.Unlock()

	if wasRunning {
		r.Reset()
	}
}

// GetInterval returns the current rest interval.
func (r *RestTimer) GetInterval() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.interval
}

// SetMaxOpacity updates the maximum overlay opacity.
// Requirements 6.6: new opacity used for future displays.
func (r *RestTimer) SetMaxOpacity(opacity float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}
	r.maxOpacity = opacity
}

// GetMaxOpacity returns the current maximum overlay opacity.
func (r *RestTimer) GetMaxOpacity() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.maxOpacity
}

// SetDisplayTime updates how long the overlay remains visible.
func (r *RestTimer) SetDisplayTime(d time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.displayTime = d
}

// GetDisplayTime returns the current overlay display duration.
func (r *RestTimer) GetDisplayTime() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.displayTime
}

// SetEnabled enables or disables the rest reminder.
func (r *RestTimer) SetEnabled(enabled bool) {
	r.mu.Lock()
	r.enabled = enabled
	r.mu.Unlock()

	if enabled {
		r.Start()
	} else {
		r.Stop()
	}
}

// IsEnabled returns whether the rest reminder is enabled.
func (r *RestTimer) IsEnabled() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.enabled
}

// Start begins (or restarts) the rest interval countdown.
// Requirements 6.1: timer fires after the configured interval.
func (r *RestTimer) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stopLocked()

	r.running = true
	interval := r.interval
	cb := r.onTrigger

	var t *time.Timer
	if r.newTimer != nil {
		t = r.newTimer(interval)
	} else {
		t = time.NewTimer(interval)
	}
	r.timer = t

	go func() {
		<-t.C
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		if cb != nil {
			cb()
		}
	}()
}

// Stop cancels the running timer without triggering.
func (r *RestTimer) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopLocked()
}

// stopLocked stops the timer. Caller must hold the write lock.
func (r *RestTimer) stopLocked() {
	r.running = false
	if r.timer != nil {
		r.timer.Stop()
		// Drain the channel if it fired between Stop and now.
		select {
		case <-r.timer.C:
		default:
		}
		r.timer = nil
	}
}

// Reset stops the current timer and starts a fresh countdown.
// Requirements 6.4: after overlay dismiss, timer restarts with same interval.
func (r *RestTimer) Reset() {
	r.Stop()
	r.Start()
}

// IsRunning returns whether the timer is currently counting down.
func (r *RestTimer) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}
