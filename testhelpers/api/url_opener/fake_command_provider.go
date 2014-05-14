package url_opener

import (
	"github.com/cloudfoundry/cli/cf/api/application_opener"
)

// implements RunnableCommand
type fakeCommand struct {
	Name         string
	Args         []string
	RunWasCalled bool
}

func (fake *fakeCommand) Run() error {
	fake.RunWasCalled = true
	return nil
}

func (fake *fakeCommand) Output() ([]byte, error) {
	return []byte{}, nil
}

// returns RunnableCommands, holding onto a reference to each
type FakeCommandProvider struct {
	CommandsProvided []*fakeCommand
}

func (fake *FakeCommandProvider) NewCommand(name string, args ...string) application_opener.RunnableCommand {
	cmd := &fakeCommand{Name: name, Args: args}
	fake.CommandsProvided = append(fake.CommandsProvided, cmd)

	return cmd
}
