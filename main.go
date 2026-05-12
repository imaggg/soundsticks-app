package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/imaggg/soundsticks/internal/api"
	"github.com/imaggg/soundsticks/internal/autostart"
	"github.com/imaggg/soundsticks/internal/cert"
	"github.com/imaggg/soundsticks/internal/config"
	"github.com/imaggg/soundsticks/internal/discovery"
	webview "github.com/webview/webview_go"
)

//go:embed ui
var uiFiles embed.FS

var (
	wv  webview.WebView
	mu  sync.Mutex
	cli *api.Client
	cfg *config.Config

	pollGen       int
	keepAliveGen  int
	discoveryGen  int
	discoveryCxl  context.CancelFunc
	momentPlaying bool

	popoverMu      sync.Mutex
	popoverVisible bool
)

// ── Entry point ───────────────────────────────────────────────────────────────

func main() {
	runtime.LockOSThread()

	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Printf("config: %v", err)
		cfg = &config.Config{}
	}

	initTray()
	initLanguage(cfg.Language)
	initTheme(cfg.Theme)
	setAutostartEnabled(autostart.IsEnabled())
	setLedKeepAliveObjC(cfg.LedKeepAlive)

	wv = webview.New(false)
	defer wv.Destroy()

	wv.SetTitle("SoundSticks")
	wv.SetSize(400, 460, webview.HintFixed)

	// Reconfigure the webview's NSWindow as a borderless popover anchored to the tray icon.
	attachPopoverWindow(wv.Window())

	port := serveUI()
	bindBridge()

	// Inject page-ready signal before any page script runs.
	// sendCmd is already registered via Bind, so it's available at load time.
	wv.Init(`window.addEventListener('load',function(){window.sendCmd('pageReady','{}');});`)

	wv.Navigate(fmt.Sprintf("http://127.0.0.1:%d/", port))

	wv.Run()
}

// ── UI server ─────────────────────────────────────────────────────────────────

func serveUI() int {
	sub, err := fs.Sub(uiFiles, "ui")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	go http.Serve(ln, mux)
	return ln.Addr().(*net.TCPAddr).Port
}

// ── JS bridge ─────────────────────────────────────────────────────────────────

func bindBridge() {
	wv.Bind("sendCmd", func(command, payload string) {
		var raw map[string]json.RawMessage
		json.Unmarshal([]byte(payload), &raw)
		go handleCmd(command, raw)
	})
}

func handleCmd(command string, payload map[string]json.RawMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("handleCmd %s panic: %v", command, r)
		}
	}()
	switch command {

	case "pageReady":
		if cert.IsSetupDone() {
			eval("showScreen('screen-searching')")
			go startDiscovery()
		} else {
			eval("showScreen('screen-setup')")
		}

	// ── Setup: APK / XAPK drop ───────────────────────────────────────────────
	case "setupImportApk":
		path := jsonStr(payload["path"])
		if path == "" {
			evalStatus(false, "Drag & drop requires the full path. Use the Browse button.")
			return
		}
		var err error
		if strings.HasSuffix(strings.ToLower(path), ".xapk") {
			err = cert.ExtractFromXAPK(path)
		} else {
			err = cert.ExtractFromAPK(path)
		}
		if err != nil {
			evalStatus(false, "Error: "+err.Error())
		} else {
			eval("window.setupDone()")
			go startDiscovery()
		}

	case "setupBrowseApk":
		go func() {
			path, err := nativeFileDialog("Select HK APK or XAPK")
			if err != nil || path == "" {
				return
			}
			var extractErr error
			if strings.HasSuffix(strings.ToLower(path), ".xapk") {
				extractErr = cert.ExtractFromXAPK(path)
			} else {
				extractErr = cert.ExtractFromAPK(path)
			}
			if extractErr != nil {
				evalStatus(false, "Error: "+extractErr.Error())
			} else {
				eval("window.setupDone()")
				go startDiscovery()
			}
		}()

	// ── Light controls ────────────────────────────────────────────────────────
	case "setLightInfo":
		c := getClient()
		if c == nil {
			return
		}
		if err := c.SetLightInfo(payload); err != nil {
			log.Printf("setLightInfo: %v", err)
		}

	// ── EQ controls ───────────────────────────────────────────────────────────
	case "setActiveEQ":
		c := getClient()
		if c == nil {
			return
		}
		var eqID int
		json.Unmarshal(payload["eq_id"], &eqID)
		if err := c.SetActiveEQ(eqID); err != nil {
			log.Printf("setActiveEQ: %v", err)
		}

	case "saveCustomPreset":
		var slot int
		var name string
		var gain []float64
		json.Unmarshal(payload["slot"], &slot)
		json.Unmarshal(payload["name"], &name)
		json.Unmarshal(payload["gain"], &gain)
		if slot >= 0 && slot < 3 && len(gain) == 7 {
			mu.Lock()
			cfg.CustomPresets[slot].Name = name
			cfg.CustomPresets[slot].Gain = gain
			err := cfg.Save()
			mu.Unlock()
			if err != nil {
				log.Printf("saveCustomPreset: %v", err)
			}
		}

	case "setCustomEQ":
		c := getClient()
		if c == nil {
			return
		}
		var fsArr []float64
		var gain []float64
		json.Unmarshal(payload["fs"], &fsArr)
		json.Unmarshal(payload["gain"], &gain)
		if err := c.SetCustomEQ(fsArr, gain); err != nil {
			log.Printf("setCustomEQ: %v", err)
		}

	case "setWindowHeight":
		var h int
		json.Unmarshal(payload["height"], &h)
		if h > 50 {
			setWindowHeightGo(h)
		}

	case "reconnect":
		eval("showScreen('screen-searching')")
		mu.Lock()
		cli = nil
		mu.Unlock()
		go startDiscovery()

	// ── Moment controls ───────────────────────────────────────────────────────
	case "momentOn":
		c := getClient()
		if c == nil {
			return
		}
		var soundscapeID, sleepTimer int
		json.Unmarshal(payload["soundscape_id"], &soundscapeID)
		json.Unmarshal(payload["sleep_timer"], &sleepTimer)
		mcfg := api.MomentConfig{SoundscapeID: soundscapeID, Volume: 70, SleepTimer: sleepTimer}
		if err := c.SetSmartButtonConfig(mcfg); err != nil {
			log.Printf("setSmartButtonConfig: %v", err)
		}
		if err := c.ControlSoundscapeV2(6, soundscapeID); err != nil {
			log.Printf("controlSoundscapeV2 on: %v", err)
		}
		mu.Lock()
		momentPlaying = true
		mu.Unlock()

	case "momentOff":
		c := getClient()
		if c == nil {
			return
		}
		var soundscapeID int
		json.Unmarshal(payload["soundscape_id"], &soundscapeID)
		if err := c.ControlSoundscapeV2(7, soundscapeID); err != nil {
			log.Printf("controlSoundscapeV2 off: %v", err)
		}
		mu.Lock()
		momentPlaying = false
		mu.Unlock()

	case "momentSwitch":
		// Atomically switch mode while playing: stop old → reconfigure → start new.
		c := getClient()
		if c == nil {
			return
		}
		var prevID, newID, sleepTimer int
		json.Unmarshal(payload["prev_id"], &prevID)
		json.Unmarshal(payload["soundscape_id"], &newID)
		json.Unmarshal(payload["sleep_timer"], &sleepTimer)
		c.ControlSoundscapeV2(7, prevID)
		mcfg := api.MomentConfig{SoundscapeID: newID, Volume: 70, SleepTimer: sleepTimer}
		if err := c.SetSmartButtonConfig(mcfg); err != nil {
			log.Printf("setSmartButtonConfig: %v", err)
		}
		c.ControlSoundscapeV2(6, newID)
		mu.Lock()
		momentPlaying = true
		mu.Unlock()

	case "momentConfig":
		// Update timer while playing without restarting.
		c := getClient()
		if c == nil {
			return
		}
		var soundscapeID, sleepTimer int
		json.Unmarshal(payload["soundscape_id"], &soundscapeID)
		json.Unmarshal(payload["sleep_timer"], &sleepTimer)
		mcfg := api.MomentConfig{SoundscapeID: soundscapeID, Volume: 70, SleepTimer: sleepTimer}
		if err := c.SetSmartButtonConfig(mcfg); err != nil {
			log.Printf("setSmartButtonConfig: %v", err)
		}

	case "momentSetElement":
		c := getClient()
		if c == nil {
			return
		}
		var soundscapeID, elementID, value int
		json.Unmarshal(payload["soundscape_id"], &soundscapeID)
		json.Unmarshal(payload["element_id"], &elementID)
		json.Unmarshal(payload["value"], &value)
		if err := c.SetSoundscapeElement(soundscapeID, elementID, value); err != nil {
			log.Printf("setSoundscapeElement: %v", err)
		}

	// ── Player controls ───────────────────────────────────────────────────────
	case "setVol":
		c := getClient()
		if c == nil {
			return
		}
		var vol int
		json.Unmarshal(payload["vol"], &vol)
		if err := c.SetPlayerVol(vol); err != nil {
			log.Printf("setPlayerVol: %v", err)
		}

	case "toggleMute":
		c := getClient()
		if c == nil {
			return
		}
		var mute bool
		json.Unmarshal(payload["mute"], &mute)
		if err := c.SetPlayerMute(mute); err != nil {
			log.Printf("setPlayerMute: %v", err)
		}
	}
}

// ── Discovery & connect ───────────────────────────────────────────────────────

func startDiscovery() {
	eval("showScreen('screen-searching')")
	tlsCert, err := cert.LoadTLSCert()
	if err != nil {
		log.Printf("load cert: %v", err)
		evalStatus(false, "Bad certificate — re-import your APK.")
		eval("showScreen('screen-setup')")
		return
	}
	log.Println("discovery: cert loaded OK, starting mDNS browse")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mu.Lock()
	if discoveryCxl != nil {
		discoveryCxl()
	}
	discoveryCxl = cancel
	discoveryGen++
	gen := discoveryGen
	mu.Unlock()

	devices, err := discovery.Browse(ctx)
	if err != nil {
		log.Printf("discovery: %v", err)
		return
	}

	for dev := range devices {
		mu.Lock()
		stale := gen != discoveryGen
		mu.Unlock()
		if stale {
			log.Println("discovery: superseded by newer discovery, exiting")
			return
		}

		log.Printf("discovery: found %s (%s), trying API", dev.Hostname, dev.IP)
		c := api.NewClient(tlsCert, dev.IP)
		info, err := c.GetDeviceInfo()
		if err != nil {
			log.Printf("getDeviceInfo %s: %v", dev.IP, err)
			continue
		}
		log.Printf("connected: %v", info)

		mu.Lock()
		if gen != discoveryGen {
			mu.Unlock()
			return
		}
		cli = c
		mu.Unlock()

		// name is nested: {"device_info": {"name": "..."}, "error_code": 0}
		name := "SoundSticks"
		if di, ok := info["device_info"].(map[string]interface{}); ok {
			if n, ok := di["name"].(string); ok && n != "" {
				name = n
			}
		}

		state := initialState(c, name)
		stateJSON, _ := json.Marshal(state)
		evalf("window.updateState(%s)", stateJSON)
		eval("showScreen('screen-main')")
		startPolling(c)
		mu.Lock()
		ka := cfg.LedKeepAlive
		mu.Unlock()
		startLedKeepAlive(c, ka)
		return
	}

	log.Println("discovery: timed out, no device responded")
	eval(`document.querySelector('#screen-searching p').textContent=window.t('searching.not_found')`)
}

func initialState(c *api.Client, deviceName string) map[string]interface{} {
	state := map[string]interface{}{
		"deviceName": deviceName,
		"connected":  true,
	}

	if light, err := c.GetLightInfo(); err != nil {
		log.Printf("getLightInfo: %v", err)
	} else {
		state["lights"] = map[string]interface{}{
			"enable":     int(light.Enable) == 1,
			"brightness": int(light.Brightness),
			"speed":      int(light.DynamicLevel),
			"patternId":  int(light.ActivePattern.ID),
			"colorLevel": int(light.ActivePattern.Level),
		}
	}

	if eq, err := c.GetEQList(); err != nil {
		log.Printf("getEQList: %v", err)
	} else {
		state["eq"] = buildEQState(eq)
	}

	if ps, err := c.GetPlayerStatus(); err != nil {
		log.Printf("getPlayerStatus: %v", err)
	} else {
		state["player"] = map[string]interface{}{
			"vol":  int(ps.Vol),
			"mute": int(ps.Mute) == 1,
		}
	}

	mu.Lock()
	playing := momentPlaying
	mu.Unlock()
	if sbc, err := c.GetSmartButtonConfig(); err != nil {
		log.Printf("getSmartButtonConfig: %v", err)
	} else {
		state["moment"] = map[string]interface{}{
			"soundscapeId": sbc.SoundscapeID,
			"sleepTimer":   sbc.SleepTimer,
			"enabled":      playing,
		}
	}

	if entries, err := c.GetSoundscapeV2Config(); err != nil {
		log.Printf("getSoundscapeV2Config: %v", err)
	} else {
		state["momentElements"] = buildElementVolumes(entries)
	}

	mu.Lock()
	state["customPresets"] = cfg.CustomPresets
	state["language"] = cfg.Language
	state["theme"] = cfg.Theme
	mu.Unlock()

	return state
}

// buildEQState extracts presetId and official custom gains from a getEQList response.
// Custom gains are doubled (device stores half of the UI-facing dB value).
func buildEQState(eq *api.EQListResponse) map[string]interface{} {
	m := map[string]interface{}{"presetId": int(eq.ActiveEQID)}
	for _, entry := range eq.EQList {
		if int(entry.EQID) == 0 {
			ui := make([]float64, len(entry.EQPayload.Gain))
			for i, g := range entry.EQPayload.Gain {
				if i == 0 && g < 0 {
					ui[i] = g * 12.0 / 9.0 // 125Hz negative: -9 API → -12 UI
				} else {
					ui[i] = g * 2
				}
			}
			m["customGain"] = ui
			break
		}
	}
	return m
}

func getClient() *api.Client {
	mu.Lock()
	defer mu.Unlock()
	return cli
}

func startPolling(c *api.Client) {
	mu.Lock()
	pollGen++
	gen := pollGen
	mu.Unlock()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			same := pollGen == gen && cli == c
			mu.Unlock()
			if !same {
				return
			}
			popoverMu.Lock()
			vis := popoverVisible
			popoverMu.Unlock()
			if !vis {
				continue
			}
			state := fetchCurrentState(c)
			if len(state) == 0 {
				continue
			}
			stateJSON, _ := json.Marshal(state)
			evalf("window.updateState(%s)", stateJSON)
		}
	}()
}

func fetchCurrentState(c *api.Client) map[string]interface{} {
	state := map[string]interface{}{}
	if light, err := c.GetLightInfo(); err == nil {
		state["lights"] = map[string]interface{}{
			"enable":     int(light.Enable) == 1,
			"brightness": int(light.Brightness),
			"speed":      int(light.DynamicLevel),
			"patternId":  int(light.ActivePattern.ID),
			"colorLevel": int(light.ActivePattern.Level),
		}
	}
	if eq, err := c.GetEQList(); err == nil {
		state["eq"] = buildEQState(eq)
	}
	if ps, err := c.GetPlayerStatus(); err == nil {
		state["player"] = map[string]interface{}{
			"vol":  int(ps.Vol),
			"mute": int(ps.Mute) == 1,
		}
	}
	mu.Lock()
	playing := momentPlaying
	mu.Unlock()
	if sbc, err := c.GetSmartButtonConfig(); err == nil {
		state["moment"] = map[string]interface{}{
			"soundscapeId": sbc.SoundscapeID,
			"sleepTimer":   sbc.SleepTimer,
			"enabled":      playing,
		}
	}
	if entries, err := c.GetSoundscapeV2Config(); err == nil {
		state["momentElements"] = buildElementVolumes(entries)
	}
	return state
}

// buildElementVolumes converts device element list (id reversed) to sliderIndex order.
// Device: element_id 2=first slider, 1=second, 0=third.
func buildElementVolumes(entries []api.SoundscapeV2Entry) map[int][]int {
	m := map[int][]int{}
	for _, entry := range entries {
		modeID := int(entry.SoundscapeID)
		vols := make([]int, 3)
		for _, el := range entry.Elements {
			sliderIdx := 2 - int(el.ID)
			if sliderIdx >= 0 && sliderIdx < 3 {
				vols[sliderIdx] = int(el.Value)
			}
		}
		m[modeID] = vols
	}
	return m
}

// ── LED keepalive ─────────────────────────────────────────────────────────────

func startLedKeepAlive(c *api.Client, totalMinutes int) {
	mu.Lock()
	keepAliveGen++
	gen := keepAliveGen
	mu.Unlock()

	if totalMinutes == 0 {
		return
	}

	go func() {
		deadline := time.Now().Add(time.Duration(totalMinutes) * time.Minute)
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if time.Now().After(deadline) {
				return
			}
			mu.Lock()
			same := keepAliveGen == gen && cli == c
			mu.Unlock()
			if !same {
				return
			}
			nudgeLedBrightness(c)
		}
	}()
}

func nudgeLedBrightness(c *api.Client) {
	light, err := c.GetLightInfo()
	if err != nil || int(light.Enable) == 0 {
		return
	}
	b := int(light.Brightness)
	nudge := b + 1
	if nudge > 100 {
		nudge = b - 1
	}
	set := func(v int) {
		raw, _ := json.Marshal(fmt.Sprintf("%d", v))
		c.SetLightInfo(map[string]json.RawMessage{"brightness": raw})
	}
	set(nudge)
	time.Sleep(200 * time.Millisecond)
	set(b)
}

// ── Webview JS helpers ────────────────────────────────────────────────────────

func eval(js string) {
	wv.Dispatch(func() { wv.Eval(js) })
}

func evalf(format string, args ...interface{}) {
	eval(fmt.Sprintf(format, args...))
}

func evalStatus(ok bool, msg string) {
	okJS := "false"
	if ok {
		okJS = "true"
	}
	evalf("setupShowStatus(%s,%q)", okJS, msg)
}

func jsonStr(v json.RawMessage) string {
	var s string
	json.Unmarshal(v, &s)
	return s
}
