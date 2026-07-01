//go:build windows

package infrastructure

import (
	"fmt"
	"io"
	"log"

	"github.com/UserExistsError/conpty"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
)

type PTYManager struct {
	shellFinder application.ShellFinder
}

var _ application.PTYProvider = (*PTYManager)(nil)

func NewPTYManager(shellFinder application.ShellFinder) *PTYManager {
	return &PTYManager{shellFinder: shellFinder}
}

func (pm *PTYManager) Stop(inst application.PTYInstance) error {
	return inst.Close()
}

func (pm *PTYManager) Spawn() (application.PTYInstance, error) {
	if !conpty.IsConPtyAvailable() {
		return nil, fmt.Errorf("conpty is not available on this version of Windows (requires Windows 10 1809+)")
	}

	shell, err := pm.shellFinder.FindShell()
	if err != nil {
		return nil, err
	}

	cpty, err := conpty.Start(shell)
	if err != nil {
		return nil, fmt.Errorf("failed to start conpty: %w", err)
	}

	instance := &PTYInstance{
		cpty:    cpty,
		output:  make(chan []byte, 256),
		done:    make(chan struct{}),
	}

	go instance.readLoop()

	return instance, nil
}

type PTYInstance struct {
	cpty    *conpty.ConPty
	output  chan []byte
	done    chan struct{}
}

func (p *PTYInstance) Write(data []byte) error {
	_, err := p.cpty.Write(data)
	return err
}

func (p *PTYInstance) Output() <-chan []byte {
	return p.output
}

func (p *PTYInstance) Close() error {
	return p.cpty.Close()
}

func (p *PTYInstance) readLoop() {
	defer close(p.done)
	defer close(p.output)

	buf := make([]byte, 4096)
	for {
		n, err := p.cpty.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("pty read error: %v", err)
			}
			return
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		select {
		case p.output <- data:
		default:
			log.Printf("pty output buffer full, dropping data")
		}
	}
}
