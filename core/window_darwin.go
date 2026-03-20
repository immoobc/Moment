//go:build darwin

package core

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setWindowLevel(int level) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSWindow *keyWindow = [NSApp keyWindow];
        if (keyWindow != nil) {
            if (level == 0) {
                [keyWindow setLevel:NSFloatingWindowLevel];
            } else {
                [keyWindow setLevel:NSNormalWindowLevel];
            }
        }
    });
}
*/
import "C"

import "fyne.io/fyne/v2"

// applyWindowLevel uses macOS NSWindow.setLevel to control the window z-order.
func applyWindowLevel(win fyne.Window, level WindowLevel) {
	C.setWindowLevel(C.int(level))
}

// RemoveTitleBar is handled by the splash window on macOS; this is a no-op fallback.
func RemoveTitleBar() {}

// moveWindowTo is a no-op on macOS; splash windows handle drag natively.
func moveWindowTo(x, y float32) {}

// GetCursorScreenPos is a no-op stub on macOS.
func GetCursorScreenPos() (int32, int32) { return 0, 0 }

// GetWindowScreenRect is a no-op stub on macOS.
func GetWindowScreenRect() (int32, int32, int32, int32) { return 0, 0, 0, 0 }
