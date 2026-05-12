APP     := SoundSticks
BINARY  := soundsticks
BUNDLE  := dist/$(APP).app
WINEXE  := dist/$(APP).exe
LINBIN  := dist/$(APP)

.PHONY: mac win linux icon-mac icon-win clean

# ── macOS: binary + .app bundle ───────────────────────────────────────────────
mac: dist icon-mac
	go build -o $(BUNDLE)/Contents/MacOS/$(APP) .
	cp Info.plist $(BUNDLE)/Contents/
	@if [ -f assets/icon.icns ]; then \
		cp assets/icon.icns $(BUNDLE)/Contents/Resources/AppIcon.icns; \
	else \
		echo "⚠  No assets/icon.icns — run 'make icon-mac' first (needs rsvg-convert)"; \
	fi
	@echo "✅ $(BUNDLE)"

# ── Windows: .exe with embedded icon ──────────────────────────────────────────
win: dist icon-win
	@which go-winres > /dev/null 2>&1 || \
		(echo "Install go-winres: go install github.com/tc-hib/go-winres@latest" && exit 1)
	go-winres make --in assets/winres/winres.json
	GOOS=windows GOARCH=amd64 go build -o $(WINEXE) .
	rm -f resource_windows_*.syso
	@echo "✅ $(WINEXE)"

# ── Linux: plain binary ───────────────────────────────────────────────────────
linux: dist
	GOOS=linux GOARCH=amd64 go build -o $(LINBIN) .
	@echo "✅ $(LINBIN)"
	@echo "ℹ  Place assets/icon.png at /usr/share/icons/hicolor/256x256/apps/soundsticks.png"
	@echo "   for the desktop icon to appear (used by autostart .desktop file)."

# ── Icon: SVG → .icns (macOS) ─────────────────────────────────────────────────
# Requires: rsvg-convert (brew install librsvg) + iconutil (Xcode CLI tools)
icon-mac: assets/icon.icns

assets/icon.icns: assets/icon.svg
	@which rsvg-convert > /dev/null 2>&1 || \
		(echo "Install rsvg-convert: brew install librsvg" && exit 1)
	mkdir -p /tmp/soundsticks.iconset
	for sz in 16 32 64 128 256 512; do \
		rsvg-convert -w $$sz -h $$sz assets/icon.svg \
			-o /tmp/soundsticks.iconset/icon_$${sz}x$${sz}.png; \
		rsvg-convert -w $$((sz*2)) -h $$((sz*2)) assets/icon.svg \
			-o /tmp/soundsticks.iconset/icon_$${sz}x$${sz}@2x.png; \
	done
	iconutil -c icns /tmp/soundsticks.iconset -o assets/icon.icns
	rm -rf /tmp/soundsticks.iconset
	@echo "✅ assets/icon.icns"

# ── Icon: SVG → .ico (Windows) ────────────────────────────────────────────────
# Requires: rsvg-convert + ImageMagick (brew install imagemagick)
icon-win: assets/winres/icon.ico

assets/winres/icon.ico: assets/icon.svg
	@which rsvg-convert > /dev/null 2>&1 || \
		(echo "Install rsvg-convert: brew install librsvg" && exit 1)
	@which convert > /dev/null 2>&1 || \
		(echo "Install ImageMagick: brew install imagemagick" && exit 1)
	for sz in 16 32 48 64 128 256; do \
		rsvg-convert -w $$sz -h $$sz assets/icon.svg -o /tmp/icon_$${sz}.png; \
	done
	convert /tmp/icon_16.png /tmp/icon_32.png /tmp/icon_48.png \
	        /tmp/icon_64.png /tmp/icon_128.png /tmp/icon_256.png \
	        assets/winres/icon.ico
	rm -f /tmp/icon_*.png
	@echo "✅ assets/winres/icon.ico"

# ── Helpers ───────────────────────────────────────────────────────────────────
dist:
	mkdir -p dist
	mkdir -p $(BUNDLE)/Contents/MacOS
	mkdir -p $(BUNDLE)/Contents/Resources

clean:
	rm -rf dist assets/icon.icns assets/winres/icon.ico resource_windows_*.syso
