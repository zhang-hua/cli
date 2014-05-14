package application_opener

import (
	"os/exec"
)

type RunnableCommand interface {
	Run() error
	Output() ([]byte, error)
}

type CommandProvider interface {
	NewCommand(name string, args ...string) RunnableCommand
}

type RealCommandProvider struct{}

func (provider RealCommandProvider) NewCommand(name string, args ...string) RunnableCommand {
	return exec.Command(name, args...)
}
