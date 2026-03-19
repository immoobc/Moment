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
