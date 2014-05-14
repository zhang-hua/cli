package application

import (
	"github.com/cloudfoundry/cli/cf/api/application_opener"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type openAppCommand struct {
	ui     terminal.UI
	urlOpener application_opener.URLOpener
	appReq requirements.ApplicationRequirement
}

func NewOpenApp(ui terminal.UI, urlOpener application_opener.URLOpener) *openAppCommand {
	return &openAppCommand{ui: ui, urlOpener: urlOpener}
}

func (cmd *openAppCommand) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "open",
		Description: "Open app in default browser",
		Usage:       "CF_NAME open APP",
	}
}

func (cmd *openAppCommand) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) < 1 {
		cmd.ui.FailWithUsage(c)
		return
	}

	cmd.appReq = requirementsFactory.NewApplicationRequirement(c.Args()[0])

	return []requirements.Requirement{cmd.appReq}, nil
}

func (cmd *openAppCommand) Run(context *cli.Context) {
	app := cmd.appReq.GetApplication()
	cmd.urlOpener.u
}
