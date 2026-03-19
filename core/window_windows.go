//go:build windows

package core

import (
	"syscall"

	"fyne.io/fyne/v2"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
)

const (
	hwndTopMost   = ^uintptr(0) // HWND_TOPMOST = -1
	hwndNoTopMost = ^uintptr(1) // HWND_NOTOPMOST = -2
	swpNoMove     = 0x0002
	swpNoSize     = 0x0001
	swpShowWindow = 0x0040
)

// applyWindowLevel uses the Windows user32.dll SetWindowPos API to set the
// window z-order. It retrieves the foreground window handle as a best-effort
// approach since Fyne does not expose the native HWND directly.
func applyWindowLevel(win fyne.Window, level WindowLevel) {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return
	}

	var insertAfter uintptr
	switch level {
	case LevelTopMost:
		insertAfter = hwndTopMost
	default:
		insertAfter = hwndNoTopMost
	}

	procSetWindowPos.Call(
		hwnd,
		insertAfter,
		0, 0, 0, 0,
		uintptr(swpNoMove|swpNoSize|swpShowWindow),
	)
}
