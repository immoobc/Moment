//go:build !windows

package core

// EnsureSingleInstance is a no-op on non-Windows platforms.
// Returns true so the app always starts.
func EnsureSingleInstance() bool {
	return true
}
