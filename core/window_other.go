//go:build !windows && !darwin

package core

import "fyne.io/fyne/v2"

// applyWindowLevel is a no-op on unsupported platforms.
// Fyne does not provide a cross-platform window level API,
// so on Linux and other OSes this is a best-effort stub.
func applyWindowLevel(win fyne.Window, level WindowLevel) {
	// No platform-specific implementation available.
}
