//go:build darwin

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

const plistLabel = "com.imaggg.soundsticks"

const plistTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.imaggg.soundsticks</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
</dict>
</plist>
`

func plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(d, 0755); err != nil {
		return "", err
	}
	return filepath.Join(d, plistLabel+".plist"), nil
}

func IsEnabled() bool {
	p, err := plistPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

func Enable() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}
	p, err := plistPath()
	if err != nil {
		return err
	}
	return os.WriteFile(p, []byte(fmt.Sprintf(plistTmpl, exe)), 0644)
}

func Disable() error {
	p, err := plistPath()
	if err != nil {
		return err
	}
	return os.Remove(p)
}
