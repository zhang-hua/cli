package serviceplan_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
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

})
