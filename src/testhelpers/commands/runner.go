package commands

import (
	"cf/commands"
	"github.com/codegangsta/cli"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"cf/app"
)

var CommandDidPassRequirements bool

type FakeRunner struct {
	command commands.Command
	requirementFactory *testreq.FakeReqFactory
}

func (fake FakeRunner) RunCmdByName(cmdName string, context *cli.Context) (err error) {
	println("in runcmdbyname")
	defer func() {
		errMsg := recover()

		if errMsg != nil && errMsg != testterm.FailedWasCalled {
			panic(errMsg)
		}
	}()

	CommandDidPassRequirements = false

	requirements, err := fake.command.GetRequirements(fake.requirementFactory, context)
	if err != nil {
		return
	}

	for _, requirement := range requirements {
		success := requirement.Execute()
		if !success {
			return
		}
	}

	CommandDidPassRequirements = true
	fake.command.Run(context)

	return
}

func RunCommand(cmd commands.Command, context *cli.Context, reqFactory *testreq.FakeReqFactory) {
	FakeRunner{cmd, reqFactory}.RunCmdByName("command-name", context)
}

func Run(commandName string, cmd commands.Command, args []string, requirementFactory *testreq.FakeReqFactory) {
	cmdRunner := FakeRunner{cmd, requirementFactory}
	testApp, _ := app.NewApp(cmdRunner)
	testApp.Run(append([]string{"cf", commandName}, args...))
}
