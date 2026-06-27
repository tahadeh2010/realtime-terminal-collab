package infrastructure

import (
	"os"
	"os/exec"
)

type PTYInstance struct {
	cmd    *exec.Cmd
	pty    *os.File
	output chan []byte
	done   chan struct{}
}

func (p *PTYInstance) Write(data []byte) error {
	_, err := p.pty.Write(data)
	return err
}

func (p *PTYInstance) Output() <-chan []byte {
	return p.output
}

func (p *PTYInstance) Done() <-chan struct{} {
	return p.done
}

func (p *PTYInstance) Close() error {
	p.cmd.Process.Kill()
	return p.pty.Close()
}
