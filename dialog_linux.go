//go:build linux

package main

import (
	"os/exec"
	"strings"
)

func nativeFileDialog(prompt string) (string, error) {
	out, err := exec.Command("zenity", "--file-selection", "--title="+prompt).Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}
