//go:build linux

package main

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

void soundsticks_present(void* w) {
    gtk_window_present(GTK_WINDOW(w));
}
void soundsticks_hide(void* w) {
    gtk_widget_hide(GTK_WIDGET(w));
}
// Returns 1 if the window is currently visible.
int soundsticks_visible(void* w) {
    return gtk_widget_is_visible(GTK_WIDGET(w)) ? 1 : 0;
}
*/
import "C"

import (
	"os"
	"runtime"
	"unsafe"

	"github.com/getlantern/systray"
)

var gtkWindow unsafe.Pointer

func init() {
	popoverMu.Lock()
	popoverVisible = true
	popoverMu.Unlock()
}

func initTray() {
	go func() {
		runtime.LockOSThread()
		systray.Run(onSystrayReady, nil)
	}()
}

func attachPopoverWindow(p unsafe.Pointer) {
	gtkWindow = p
}

func setWindowHeightGo(h int)        {}
func initLanguage(lang string)       {}
func initTheme(theme string)         {}
func setAutostartEnabled(v bool)     {}
func setLedKeepAliveObjC(min int)    {}

func onSystrayReady() {
	systray.SetIcon(makeTrayIcon())
	systray.SetTitle("SoundSticks")
	systray.SetTooltip("SoundSticks")

	mToggle := systray.AddMenuItem("Hide Window", "Show or hide the SoundSticks window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit SoundSticks")

	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if gtkWindow == nil {
					continue
				}
				wv.Dispatch(func() {
					if C.soundsticks_visible(gtkWindow) != 0 {
						C.soundsticks_hide(gtkWindow)
						mToggle.SetTitle("Show Window")
					} else {
						C.soundsticks_present(gtkWindow)
						mToggle.SetTitle("Hide Window")
					}
				})
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}
