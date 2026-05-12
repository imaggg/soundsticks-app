//go:build windows

package main

import (
	"os/exec"
	"strings"
)

func nativeFileDialog(prompt string) (string, error) {
	script := `Add-Type -AssemblyName System.Windows.Forms; $d = New-Object System.Windows.Forms.OpenFileDialog; $d.Title = '` + prompt + `'; if ($d.ShowDialog() -eq 'OK') { $d.FileName }`
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}
