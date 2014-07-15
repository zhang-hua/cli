package serviceplan_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/serviceplan"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("service-access command", func() {
	var (
		ui                  *testterm.FakeUI
		brokerRepo          *testapi.FakeServiceBrokerRepo
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		brokerRepo = &testapi.FakeServiceBrokerRepo{}
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	runCommand := func() bool {
		cmd := NewServiceAccess(ui, testconfig.NewRepositoryWithDefaults(), brokerRepo)
		return testcmd.RunCommand(cmd, []string{}, requirementsFactory)
	}

	Describe("requirements", func() {
		It("requires the user to be logged in", func() {
			requirementsFactory.LoginSuccess = false
			Expect(runCommand()).ToNot(HavePassedRequirements())
		})
	})

	Describe("when logged in", func() {
		BeforeEach(func() {
			serviceBrokers := []models.ServiceBroker{
				{Guid: "broker1", Name: "brokername1"},
				{Guid: "broker2", Name: "brokername2"},
			}
			brokerRepo.ServiceBrokers = serviceBrokers
		})
		It("prints all of the brokers", func() {
			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings([]string{"broker: brokername1"}, []string{"broker: brokername2"}))
		})
	})
})
