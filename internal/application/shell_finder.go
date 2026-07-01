package application

type ShellFinder interface {
	FindShell() (string, error)
}
