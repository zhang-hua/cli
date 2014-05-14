package application_test

import (
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/application"
	//	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("open command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = new(testterm.FakeUI)
		requirementsFactory = new(testreq.FakeReqFactory)
	})

	runCommand := func(args ...string) {
		cmd := NewOpenApp(ui)
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
			app.Name = "sweet-app-i-wrote"
			requirementsFactory.Application = app
		})

		It("opens the app with that name when it exists", func() {
			runCommand("sweet-app-i-wrote")

			Expect(requirementsFactory.ApplicationName).To(Equal("sweet-app-i-wrote"))
			Expect(ui.FailedWithUsage).To(BeFalse())
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
