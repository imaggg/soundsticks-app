//go:build darwin

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func nativeFileDialog(prompt string) (string, error) {
	out, err := exec.Command("osascript", "-e",
		fmt.Sprintf(`POSIX path of (choose file with prompt %q)`, prompt),
	).Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}
