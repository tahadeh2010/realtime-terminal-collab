package application

type PTYProvider interface {
	Spawn() (PTYInstance, error)
	Stop(inst PTYInstance) error
}

type PTYInstance interface {
	Write(data []byte) error
	Output() <-chan []byte
	Close() error
}
