//go:build darwin

package main

/*
#cgo LDFLAGS: -framework Cocoa
#include <stdlib.h>
void setupTray(void);
void attachPopoverWindow(void *nsWindowPtr);
void setWindowHeight(int h);
void showStatusAlert(const char *text);
void showInfoPanel(const char *name, const char *firmware,
                   const char *serial, const char *mac, const char *uuid);
void showForgetConfirm(void);
void setLanguageMenuCheck(const char *lang);
void setAutostartCheck(int enabled);
void setLedKeepAliveCheck(int minutes);
void setThemeMenuCheck(const char *theme);
void setDeviceIPField(const char *ip);
void showConnectIPPrompt(void);
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/imaggg/soundsticks/internal/autostart"
	"github.com/imaggg/soundsticks/internal/cert"
)

//export trayQuit
func trayQuit() { os.Exit(0) }

//export trayReconnect
func trayReconnect() {
	eval("showScreen('screen-searching')")
	mu.Lock()
	cli = nil
	mu.Unlock()
	go startDiscovery()
}

//export trayConnectIP
func trayConnectIP(cIP *C.char) {
	ip := strings.TrimSpace(C.GoString(cIP))
	if ip == "" {
		return
	}
	mu.Lock()
	cfg.DeviceIP = ip
	cfg.Save()
	cli = nil
	mu.Unlock()
	setTrayDeviceIP(ip)
	eval("showScreen('screen-searching')")
	// startDiscovery now tries the cached IP directly before mDNS, so this
	// connects straight to the entered address.
	go startDiscovery()
}

// setTrayDeviceIP prefills the manual "Connect to IP…" dialog with the last IP.
func setTrayDeviceIP(ip string) {
	cstr := C.CString(ip)
	C.setDeviceIPField(cstr)
	C.free(unsafe.Pointer(cstr))
}

//export trayStatus
func trayStatus() {
	go func() {
		c := getClient()
		if c == nil {
			cstr := C.CString("Not connected to SoundSticks.")
			C.showStatusAlert(cstr)
			C.free(unsafe.Pointer(cstr))
			return
		}

		var lines []string

		if ps, err := c.GetPlayerStatus(); err == nil {
			status := ps.Status
			switch status {
			case "stop":
				status = "Stopped"
			case "play":
				status = "Playing"
			}
			muteStr := "No"
			if int(ps.Mute) != 0 {
				muteStr = "Yes"
			}
			eqNames := map[int]string{0: "Custom", 1: "Signature", 2: "Vocal", 3: "Energetic", 4: "Chill"}
			eqName, ok := eqNames[int(ps.EQ)]
			if !ok {
				eqName = fmt.Sprintf("ID %d", ps.EQ)
			}
			lines = append(lines,
				fmt.Sprintf("Playback:    %s", status),
				fmt.Sprintf("Volume:      %d", ps.Vol),
				fmt.Sprintf("Muted:       %s", muteStr),
				fmt.Sprintf("EQ:          %s", eqName),
			)
		}

		if light, err := c.GetLightInfo(); err == nil {
			onOff := "OFF"
			if int(light.Enable) == 1 {
				onOff = "ON"
			}
			lines = append(lines, "", fmt.Sprintf("Lights:      %s", onOff))
			if int(light.Enable) == 1 {
				patNames := map[int]string{
					1: "Ocean", 2: "Aurora", 3: "Blossom",
					4: "Sunrise", 5: "Fireplace", 6: "Calm", 7: "Nebula",
				}
				pat, ok := patNames[int(light.ActivePattern.ID)]
				if !ok {
					pat = fmt.Sprintf("ID %d", light.ActivePattern.ID)
				}
				lines = append(lines,
					fmt.Sprintf("Pattern:     %s", pat),
					fmt.Sprintf("Brightness:  %d%%", int(light.Brightness)),
					fmt.Sprintf("Speed:       %d", int(light.DynamicLevel)),
				)
			}
		}

		text := strings.Join(lines, "\n")
		cstr := C.CString(text)
		C.showStatusAlert(cstr)
		C.free(unsafe.Pointer(cstr))
	}()
}

//export trayInfo
func trayInfo() {
	go func() {
		c := getClient()
		if c == nil {
			cstr := C.CString("Not connected to SoundSticks.")
			C.showStatusAlert(cstr)
			C.free(unsafe.Pointer(cstr))
			return
		}

		raw, err := c.GetDeviceInfo()
		if err != nil {
			cstr := C.CString("Failed to get device info: " + err.Error())
			C.showStatusAlert(cstr)
			C.free(unsafe.Pointer(cstr))
			return
		}

		di, _ := raw["device_info"].(map[string]interface{})
		str := func(k string) string {
			if v, ok := di[k].(string); ok && v != "" {
				return v
			}
			return "—"
		}

		cname := C.CString(str("name"))
		cfirmware := C.CString(str("firmware"))
		cserial := C.CString(str("serial_number"))
		cmac := C.CString(str("mac"))
		cuuid := C.CString(str("uuid"))
		C.showInfoPanel(cname, cfirmware, cserial, cmac, cuuid)
		C.free(unsafe.Pointer(cname))
		C.free(unsafe.Pointer(cfirmware))
		C.free(unsafe.Pointer(cserial))
		C.free(unsafe.Pointer(cmac))
		C.free(unsafe.Pointer(cuuid))
	}()
}

//export trayForget
func trayForget() { C.showForgetConfirm() }

//export trayToggleAutostart
func trayToggleAutostart() {
	enabled := autostart.IsEnabled()
	var err error
	if enabled {
		err = autostart.Disable()
	} else {
		err = autostart.Enable()
	}
	if err != nil {
		log.Printf("autostart: %v", err)
		return
	}
	setAutostartEnabled(!enabled)
}

func setAutostartEnabled(v bool) {
	e := C.int(0)
	if v {
		e = C.int(1)
	}
	C.setAutostartCheck(e)
}

//export traySetLedKeepAlive
func traySetLedKeepAlive(cMinutes C.int) {
	minutes := int(cMinutes)
	mu.Lock()
	cfg.LedKeepAlive = minutes
	c := cli
	cfg.Save()
	mu.Unlock()
	C.setLedKeepAliveCheck(cMinutes)
	if c != nil {
		startLedKeepAlive(c, minutes)
	}
}

func setLedKeepAliveObjC(minutes int) {
	C.setLedKeepAliveCheck(C.int(minutes))
}

//export traySetTheme
func traySetTheme(cTheme *C.char) {
	theme := C.GoString(cTheme)
	mu.Lock()
	cfg.Theme = theme
	cfg.Save()
	mu.Unlock()
	cstr := C.CString(theme)
	C.setThemeMenuCheck(cstr)
	C.free(unsafe.Pointer(cstr))
	evalf("window.applyTheme(%q)", theme)
}

//export traySetLanguage
func traySetLanguage(cLang *C.char) {
	lang := C.GoString(cLang)
	if lang == "" {
		lang = "en"
	}
	mu.Lock()
	cfg.Language = lang
	cfg.Save()
	mu.Unlock()
	cstr := C.CString(lang)
	C.setLanguageMenuCheck(cstr)
	C.free(unsafe.Pointer(cstr))
	evalf("window.setLanguage(%q)", lang)
}

//export trayForgetConfirmed
func trayForgetConfirmed() {
	cert.Forget()
	mu.Lock()
	cli = nil
	mu.Unlock()
	eval("showScreen('screen-setup')")
}

//export trayPopoverWillShow
func trayPopoverWillShow() {
	popoverMu.Lock()
	popoverVisible = true
	popoverMu.Unlock()
	// Immediately fetch and push current device state.
	go func() {
		c := getClient()
		if c == nil {
			return
		}
		state := fetchCurrentState(c)
		if len(state) == 0 {
			return
		}
		stateJSON, _ := json.Marshal(state)
		evalf("window.updateState(%s)", stateJSON)
	}()
}

//export trayPopoverDidHide
func trayPopoverDidHide() {
	popoverMu.Lock()
	popoverVisible = false
	popoverMu.Unlock()
}

func initTray()                            { C.setupTray() }
func attachPopoverWindow(p unsafe.Pointer) { C.attachPopoverWindow(p) }
func setWindowHeightGo(h int)              { C.setWindowHeight(C.int(h)) }

func initLanguage(lang string) {
	cstr := C.CString(lang)
	C.setLanguageMenuCheck(cstr)
	C.free(unsafe.Pointer(cstr))
}

func initTheme(theme string) {
	cstr := C.CString(theme)
	C.setThemeMenuCheck(cstr)
	C.free(unsafe.Pointer(cstr))
}
