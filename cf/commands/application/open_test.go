package application_test

import (
	"github.com/cloudfoundry/cli/cf/models"
	testopener "github.com/cloudfoundry/cli/testhelpers/api/url_opener"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/application"
	//	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("open command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		urlOpener           *testopener.FakeURLOpener
	)

	BeforeEach(func() {
		ui = new(testterm.FakeUI)
		requirementsFactory = new(testreq.FakeReqFactory)
		urlOpener = new(testopener.FakeURLOpener)
	})

	runCommand := func(args ...string) {
		cmd := NewOpenApp(ui, urlOpener)
		testcmd.RunCommand(cmd, testcmd.NewContext(cmd.Metadata().Name, args), requirementsFactory)
	}

	It("fails with usage when invoked with no args", func() {
		runCommand()
		Expect(ui.FailedWithUsage).To(BeTrue())
	})

	Context("when the app does exist", func() {
		BeforeEach(func() {
			requirementsFactory.AppNotFound = false

			app := models.Application{}
			app.Name = "my-app"
			app.Routes = []models.RouteSummary{
				models.RouteSummary{
					Host: "my-app",
					Domain: models.DomainFields{
						Name: "run.pivotal.io",
					},
				},
			}
			requirementsFactory.Application = app
		})

		It("opens the app with that name when it exists", func() {
			runCommand("my-app")

			Expect(requirementsFactory.ApplicationName).To(Equal("my-app"))
			Expect(ui.FailedWithUsage).To(BeFalse())
			Expect(urlOpener.OpenURLReceived.URL).To(Equal("my-app.run.pivotal.io"))
		})
	})

	Context("when the app specified does not exist", func() {
		BeforeEach(func() {
			requirementsFactory.AppNotFound = true
		})

		It("fails requirements", func() {
			runCommand("unknown-app")

			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

})
