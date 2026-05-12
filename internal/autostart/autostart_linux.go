//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

func desktopPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(home, ".config", "autostart")
	if err := os.MkdirAll(d, 0755); err != nil {
		return "", err
	}
	return filepath.Join(d, "soundsticks.desktop"), nil
}

func IsEnabled() bool {
	p, err := desktopPath()
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
	p, err := desktopPath()
	if err != nil {
		return err
	}
	content := fmt.Sprintf(
		"[Desktop Entry]\nName=SoundSticks\nExec=%s\nIcon=soundsticks\nType=Application\nX-GNOME-Autostart-enabled=true\n",
		exe,
	)
	return os.WriteFile(p, []byte(content), 0644)
}

func Disable() error {
	p, err := desktopPath()
	if err != nil {
		return err
	}
	return os.Remove(p)
}
