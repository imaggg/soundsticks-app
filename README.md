# SoundSticks 5 Wi-Fi Remote Control

Unofficial tray app to control **Harman Kardon SoundSticks 5 WiFi** speakers.  
No cloud. No account. Pure local mTLS to the device.
Tested with app v2.5.4 and speaker firmware OneOS 3.1 (26.01.21.63.00)

> **Linux & Windows builds are included but untested.**

---

## Features

- **General** — lights on/off, brightness, speed, 7 patterns, color hue
- **Moment** — on/off, 4 modes, per-mode element sliders, sleep timer (up to 60 min) _(untested)_
- **EQ** — 5 factory presets + 3 local saved user presets
- Light / Dark / Auto theme
- mDNS auto-discovery — no IP configuration needed. The last-known IP is remembered and tried directly on every reconnect, so the app stays connected even when mDNS is flaky; **Connect to IP…** in the tray menu lets you enter the address manually if discovery ever fails.
- LED sleep override / keep-alive _(untested)_
- Autostart on login
- 4 languages: English, Ukrainian, Spanish, French

---

## How it works

SoundSticks 5 WiFi runs [Linkplay](https://www.linkplay.com/) firmware. The official **Harman Kardon One** Android app communicates with the speaker over HTTPS using mutual TLS (mTLS) — both sides authenticate with certificates. The speaker's client certificate (`alice.p12`) is bundled inside the APK, password-protected with a key derived from the app's native library.

This app reverse-engineered the password derivation, so it can extract the certificate automatically. No keys or passwords are stored in this repository — only the encrypted strings already present in the publicly distributed APK.

Once extracted, the certificate is saved locally and used for all subsequent connections. Everything stays on your local network.

## Setup

Download the **Harman Kardon One** official XAPK from [APKPure](https://apkpure.com) or a similar source, then open SoundSticks and select the file on the setup screen. The app extracts the certificate, discovers the password automatically, and connects to your speaker — no manual steps required.

---

## Building

### macOS (tested)

```bash
# Prerequisites: Xcode Command Line Tools, Go 1.21+, librsvg (for icon)
xcode-select --install
brew install librsvg

make mac
open dist/SoundSticks.app
```

This builds a proper `.app` bundle in `dist/`. Double-clicking it in Finder works too — no Terminal window, no Dock icon.

> Without `librsvg`, run `go build -o dist/SoundSticks.app/Contents/MacOS/SoundSticks . && cp Info.plist dist/SoundSticks.app/Contents/` to skip the icon step.

---

### Linux (untested)

```bash
# Prerequisites — Debian/Ubuntu
sudo apt install golang libgtk-3-dev libwebkit2gtk-4.0-dev libayatana-appindicator3-dev zenity

# Prerequisites — Fedora 40 and earlier
sudo dnf install golang gtk3-devel webkit2gtk4.0-devel libayatana-appindicator-gtk3-devel zenity

# Prerequisites — Fedora 41+ (webkit2gtk-4.0 replaced by 4.1)
sudo dnf install golang gtk3-devel webkit2gtk4.1-devel libayatana-appindicator-gtk3-devel zenity
# Create a compatibility shim so the build finds webkit2gtk-4.0:
sudo tee /usr/local/lib/pkgconfig/webkit2gtk-4.0.pc > /dev/null << 'EOF'
Name: webkit2gtk-4.0
Description: shim for webkit2gtk-4.1
Version: 2.46.0
Requires: webkit2gtk-4.1
EOF

go build -o soundsticks .
./soundsticks
```

`zenity` is only needed for the file-picker dialog on the setup screen.
`libayatana-appindicator` provides the system tray icon (GNOME requires the [AppIndicator extension](https://extensions.gnome.org/extension/615/appindicator-support/)).

---

### Windows (untested)

Requirements:

- [Go 1.21+](https://go.dev/dl/)
- [MinGW-w64](https://www.mingw-w64.org/) (for CGO — add to PATH)
- [Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) (usually pre-installed on Win10+)

```powershell
go build -ldflags "-H windowsgui" -o soundsticks.exe .
.\soundsticks.exe
```

The app shows a window and a system-tray icon. Use **Show Window** from the tray to restore if minimized.

To embed a proper icon and manifest, install [go-winres](https://github.com/tc-hib/go-winres) and run `go-winres make` before building (see `assets/winres/`).

---

## Project structure

```
main.go              — app entry, JS bridge, polling, device logic
tray_darwin.go/.m    — NSStatusBar tray + popover (Cocoa/CGO)
tray_linux.go        — systray icon via getlantern/systray + GTK CGO
tray_windows.go      — systray icon via getlantern/systray
dialog_*.go          — native file-picker per platform
internal/api/        — typed mTLS client (FlexInt, eq, moment, light)
internal/cert/       — APK/XAPK extraction + .p12 import
internal/config/     — JSON config (~/.config/soundsticks or %APPDATA%)
internal/autostart/  — LaunchAgent / registry / .desktop
internal/discovery/  — mDNS browse (zeroconf)
ui/                  — single-page HTML/CSS/JS, embedded in binary
```

---

## Built with

This app was fully designed and implemented in collaboration with [Claude](https://claude.ai) (Anthropic) — architecture, API reverse-engineering, cipher reconstruction, UI, and platform integrations.

---

## License

MIT
