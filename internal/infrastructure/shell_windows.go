//go:build windows

package infrastructure

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type windowsShellFinder struct{}

func newPlatformShellFinder() *windowsShellFinder {
	return &windowsShellFinder{}
}

func (f *windowsShellFinder) FindShell() (string, error) {
	if comspec := os.Getenv("COMSPEC"); comspec != "" {
		if _, err := exec.LookPath(comspec); err == nil {
			return comspec, nil
		}
	}

	if path, err := exec.LookPath("powershell.exe"); err == nil {
		return path, nil
	}

	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}
	cmdPath := filepath.Join(systemRoot, "System32", "cmd.exe")
	if _, err := exec.LookPath(cmdPath); err == nil {
		return cmdPath, nil
	}

	return "", fmt.Errorf("no shell found")
}
