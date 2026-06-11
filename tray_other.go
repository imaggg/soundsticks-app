//go:build !darwin && !linux && !windows

package main

import "unsafe"

func init() {
	// No tray popover on non-darwin: window is always visible, poll continuously.
	popoverMu.Lock()
	popoverVisible = true
	popoverMu.Unlock()
}

func initTray()                            {}
func attachPopoverWindow(p unsafe.Pointer) {}
func setWindowHeightGo(h int)              {}
func initLanguage(lang string)             {}
func initTheme(theme string)               {}
func setAutostartEnabled(v bool)           {}
func setLedKeepAliveObjC(minutes int)      {}
func setTrayDeviceIP(ip string)            {}
