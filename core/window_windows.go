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
	dwmapi                       = syscall.NewLazyDLL("dwmapi.dll")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowLong            = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLong            = user32.NewProc("SetWindowLongPtrW")
	procMoveWindow               = user32.NewProc("MoveWindow")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procGetClientRect            = user32.NewProc("GetClientRect")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetCurrentProcessId      = kernel32.NewProc("GetCurrentProcessId")
	procDwmSetWindowAttribute    = dwmapi.NewProc("DwmSetWindowAttribute")
	procGetCursorPos             = user32.NewProc("GetCursorPos")
)

const (
	hwndTopMost     = ^uintptr(0) // HWND_TOPMOST  = -1
	hwndNoTopMost   = ^uintptr(1) // HWND_NOTOPMOST = -2
	swpNoMove       = 0x0002
	swpNoSize       = 0x0001
	swpShowWindow   = 0x0040
	swpFrameChanged = 0x0020

	gwlStyle     uintptr = ^uintptr(15) // GWL_STYLE = -16
	wsCaption    uintptr = 0x00C00000
	wsThickFrame uintptr = 0x00040000
	wsSysMenu    uintptr = 0x00080000

	// DWM attribute for window corner preference (Windows 11+)
	dwmwaWindowCornerPreference = 33
	dwmwcpRound                 = 2
)

var cachedHWND uintptr

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
					return 0
				}
			}
			return 1
		}),
		0,
	)
	return cachedHWND
}

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

// RemoveTitleBar removes the window chrome and applies DWM round corners.
//
// The key insight: after stripping WS_CAPTION | WS_THICKFRAME the window's
// non-client area shrinks, but the overall window rect stays the same.
// That means the client area (where Fyne renders) grows to fill the whole
// window rect — so content is NOT clipped. We just need to trigger a
// frame recalculation and then ask DWM for rounded corners.
func RemoveTitleBar() {
	hwnd := findOwnHWND()
	if hwnd == 0 {
		return
	}

	// Strip caption + thick frame + sys menu
	style, _, _ := procGetWindowLong.Call(hwnd, gwlStyle)
	style &^= wsCaption | wsThickFrame | wsSysMenu
	procSetWindowLong.Call(hwnd, gwlStyle, style)

	// Force Windows to recalculate the frame
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0,
		uintptr(swpNoMove|swpNoSize|swpShowWindow|swpFrameChanged))

	// Windows 11: request rounded corners via DWM (silently fails on Win10)
	pref := uint32(dwmwcpRound)
	procDwmSetWindowAttribute.Call(hwnd,
		uintptr(dwmwaWindowCornerPreference),
		uintptr(unsafe.Pointer(&pref)),
		uintptr(unsafe.Sizeof(pref)))
}

// RefreshRoundRegion is a no-op; DWM handles round corners automatically.
func RefreshRoundRegion() {}

type screenPoint struct{ X, Y int32 }

func GetCursorScreenPos() (int32, int32) {
	var pt screenPoint
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

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
