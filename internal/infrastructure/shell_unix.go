//go:build !windows

package infrastructure

import (
	"fmt"
	"os"
	"os/exec"
)

type unixShellFinder struct{}

func newPlatformShellFinder() *unixShellFinder {
	return &unixShellFinder{}
}

func (f *unixShellFinder) FindShell() (string, error) {
	if shell := os.Getenv("SHELL"); shell != "" {
		if _, err := exec.LookPath(shell); err == nil {
			return shell, nil
		}
	}

	fallback := []string{"bash", "zsh", "sh"}
	for _, name := range fallback {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no shell found")
}
