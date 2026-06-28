package infrastructure

import (
	"io"
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
)

type PTYManager struct{}

var _ application.PTYProvider = (*PTYManager)(nil)

func NewPTYManager() *PTYManager {
	return &PTYManager{}
}

func (pm *PTYManager) Stop(inst application.PTYInstance) error {
	return inst.Close()
}

func (pm *PTYManager) Spawn() (application.PTYInstance, error) {
	cmd := exec.Command("bash")
	cmd.Env = append(cmd.Env, "TERM=xterm")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	instance := &PTYInstance{
		cmd:    cmd,
		pty:    ptmx,
		output: make(chan []byte, 256),
		done:   make(chan struct{}),
	}

	go instance.readLoop()

	return instance, nil
}

func (p *PTYInstance) readLoop() {
	defer close(p.done)
	defer close(p.output)

	buf := make([]byte, 4096)
	for {
		n, err := p.pty.Read(buf)
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
