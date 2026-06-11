//go:build windows

package main

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/getlantern/systray"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procShowWindow  = user32.NewProc("ShowWindow")
	procSetFgWindow = user32.NewProc("SetForegroundWindow")
)

const swRestore = 9

var winHandle uintptr

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
	winHandle = uintptr(p)
}

func setWindowHeightGo(h int) {}
func initLanguage(lang string) {}
func initTheme(theme string)   {}
func setAutostartEnabled(v bool) {}
func setLedKeepAliveObjC(minutes int) {}
func setTrayDeviceIP(ip string) {}

func onSystrayReady() {
	systray.SetIcon(makeTrayIcon())
	systray.SetTooltip("SoundSticks")

	mShow := systray.AddMenuItem("Show Window", "Bring SoundSticks to front")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit SoundSticks")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if wv != nil {
					wv.Dispatch(func() {
						procShowWindow.Call(winHandle, swRestore)
						procSetFgWindow.Call(winHandle)
					})
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

