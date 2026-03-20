//go:build windows

package core

import (
	"os"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowLong            = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLong            = user32.NewProc("SetWindowLongPtrW")
	procMoveWindow               = user32.NewProc("MoveWindow")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetCurrentProcessId      = kernel32.NewProc("GetCurrentProcessId")
)

const (
	hwndTopMost     = ^uintptr(0) // HWND_TOPMOST = -1
	hwndNoTopMost   = ^uintptr(1) // HWND_NOTOPMOST = -2
	swpNoMove       = 0x0002
	swpNoSize       = 0x0001
	swpShowWindow   = 0x0040
	swpFrameChanged = 0x0020

	gwlStyle     uintptr = ^uintptr(15) // GWL_STYLE = -16
	wsCaption    uintptr = 0x00C00000
	wsThickFrame uintptr = 0x00040000
)

// cachedHWND stores the main window handle once found.
var cachedHWND uintptr

// findOwnHWND finds the first visible window belonging to our process.
func findOwnHWND() uintptr {
	if cachedHWND != 0 {
		return cachedHWND
	}
	pid := uint32(os.Getpid())
	procEnumWindows.Call(
		syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
			var windowPid uint32
			procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&windowPid)))
			if windowPid == pid {
				visible, _, _ := procIsWindowVisible.Call(hwnd)
				if visible != 0 {
					cachedHWND = hwnd
					return 0 // stop enumeration
				}
			}
			return 1 // continue
		}),
		0,
	)
	return cachedHWND
}

// applyWindowLevel sets the window z-order via SetWindowPos.
func applyWindowLevel(_ fyne.Window, level WindowLevel) {
	hwnd := findOwnHWND()
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
	procSetWindowPos.Call(hwnd, insertAfter, 0, 0, 0, 0,
		uintptr(swpNoMove|swpNoSize|swpShowWindow))
}

// RemoveTitleBar strips the title bar and thick frame from our window.
func RemoveTitleBar() {
	hwnd := findOwnHWND()
	if hwnd == 0 {
		return
	}
	style, _, _ := procGetWindowLong.Call(hwnd, gwlStyle)
	style &^= wsCaption | wsThickFrame
	procSetWindowLong.Call(hwnd, gwlStyle, style)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0,
		uintptr(swpNoMove|swpNoSize|swpShowWindow|swpFrameChanged))
}

var procGetCursorPos = user32.NewProc("GetCursorPos")

type screenPoint struct{ X, Y int32 }

// GetCursorScreenPos returns the current mouse cursor position in screen coordinates.
func GetCursorScreenPos() (int32, int32) {
	var pt screenPoint
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

// GetWindowScreenRect returns the window's current screen rectangle (left, top, width, height).
func GetWindowScreenRect() (int32, int32, int32, int32) {
	hwnd := findOwnHWND()
	if hwnd == 0 {
		return 0, 0, 0, 0
	}
	type rect struct{ Left, Top, Right, Bottom int32 }
	var r rect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	return r.Left, r.Top, r.Right - r.Left, r.Bottom - r.Top
}

// moveWindowTo moves our window to the given screen coordinates.
func moveWindowTo(x, y float32) {
	hwnd := findOwnHWND()
	if hwnd == 0 {
		return
	}
	type rect struct{ Left, Top, Right, Bottom int32 }
	var r rect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	w := r.Right - r.Left
	h := r.Bottom - r.Top
	procMoveWindow.Call(hwnd, uintptr(int32(x)), uintptr(int32(y)), uintptr(w), uintptr(h), 1)
}
