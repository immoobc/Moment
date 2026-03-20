//go:build !windows && !darwin

package core

import "fyne.io/fyne/v2"

// applyWindowLevel is a no-op on unsupported platforms.
// Fyne does not provide a cross-platform window level API,
// so on Linux and other OSes this is a best-effort stub.
func applyWindowLevel(win fyne.Window, level WindowLevel) {
	// No platform-specific implementation available.
}

// RemoveTitleBar is a no-op on unsupported platforms.
func RemoveTitleBar() {}

// moveWindowTo is a no-op on unsupported platforms.
func moveWindowTo(x, y float32) {}

// GetCursorScreenPos is a no-op stub on unsupported platforms.
func GetCursorScreenPos() (int32, int32) { return 0, 0 }

// GetWindowScreenRect is a no-op stub on unsupported platforms.
func GetWindowScreenRect() (int32, int32, int32, int32) { return 0, 0, 0, 0 }

// RefreshRoundRegion is a no-op on unsupported platforms.
func RefreshRoundRegion() {}
