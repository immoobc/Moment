//go:build windows

package core

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	procCreateMutexW     = kernel32.NewProc("CreateMutexW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procSetForegroundWin = user32.NewProc("SetForegroundWindow")
)

const (
	errorAlreadyExists = 183
	swShow             = 5
)

const mutexName = "Global\\MomentDesktopClock_SingleInstance"

// EnsureSingleInstance tries to create a named mutex.
// If another instance already holds it, it finds that instance's window,
// brings it to the foreground, and returns false (caller should exit).
// If this is the first instance, returns true (caller should continue).
func EnsureSingleInstance() bool {
	namePtr, _ := syscall.UTF16PtrFromString(mutexName)
	_, _, err := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(namePtr)))
	if err != nil && err.(syscall.Errno) == errorAlreadyExists {
		// Another instance is running — find its window and bring to front.
		bringExistingToFront()
		return false
	}
	return true
}

// bringExistingToFront enumerates all windows, finds one belonging to
// another moment.exe process, and brings it to the foreground.
func bringExistingToFront() {
	myPid := uint32(os.Getpid())

	// Find moment.exe processes via window enumeration
	procEnumWindows.Call(
		syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
			var pid uint32
			procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
			if pid == myPid || pid == 0 {
				return 1 // skip our own process
			}
			visible, _, _ := procIsWindowVisible.Call(hwnd)
			if visible == 0 {
				return 1
			}
			// Check if this window belongs to a moment.exe process
			if isMomentProcess(pid) {
				procShowWindow.Call(hwnd, swShow)
				procSetForegroundWin.Call(hwnd)
				// Also set topmost briefly to ensure visibility
				procSetWindowPos.Call(hwnd, hwndTopMost, 0, 0, 0, 0,
					uintptr(swpNoMove|swpNoSize|swpShowWindow))
				return 0 // stop
			}
			return 1
		}),
		0,
	)
}

// isMomentProcess checks if the given PID is a moment.exe process
// by comparing executable names.
func isMomentProcess(pid uint32) bool {
	const processQueryLimitedInfo = 0x1000
	handle, err := syscall.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(handle)

	var buf [syscall.MAX_PATH]uint16
	size := uint32(len(buf))
	ret, _, _ := procQueryFullProcessImageNameW.Call(
		uintptr(handle), 0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return false
	}
	name := syscall.UTF16ToString(buf[:size])
	// Check if the exe name ends with moment.exe
	return len(name) >= 10 && name[len(name)-10:] == "moment.exe"
}

var procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
