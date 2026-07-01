package infrastructure

import "github.com/tahadeh2010/realtime-terminal-collab/internal/application"

func NewShellFinder() application.ShellFinder {
	return newPlatformShellFinder()
}
