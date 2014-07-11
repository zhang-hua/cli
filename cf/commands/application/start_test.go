package application_test

import (
	"os"
	"time"

	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/models"
	clock "github.com/cloudfoundry/cli/clock/fakes"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testlogs "github.com/cloudfoundry/cli/testhelpers/logs"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	"github.com/cloudfoundry/loggregatorlib/logmessage"

	. "github.com/cloudfoundry/cli/cf/commands/application"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("start command", func() {
	var (
		ui                        *testterm.FakeUI
		cmd                       *Start
		mockClock                 = &clock.FakeClock{}
		app                       = models.Application{}
		defaultInstanceReponses   = [][]models.AppInstanceFields{}
		defaultInstanceErrorCodes = []string{"", ""}
		requirementsFactory       *testreq.FakeReqFactory
		configRepo                configuration.ReadWriter
		appRepo                   *testapi.FakeApplicationRepository
		appDisplayer              *testcmd.FakeAppDisplayer
		appInstancesRepo          *testapi.FakeAppInstancesRepo
		logRepo                   *testapi.FakeLogsRepository
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		appDisplayer = &testcmd.FakeAppDisplayer{}
		appRepo = &testapi.FakeApplicationRepository{}
		appInstancesRepo = &testapi.FakeAppInstancesRepo{}
		logRepo = &testapi.FakeLogsRepository{}

		app = models.Application{}
		app.Name = "my-app"
		app.Guid = "my-app-guid"
		app.InstanceCount = 2

		domain := models.DomainFields{}
		domain.Name = "example.com"

		route := models.RouteSummary{}
		route.Host = "my-app"
		route.Domain = domain

		app.Routes = []models.RouteSummary{route}

		starting := models.AppInstanceFields{State: models.InstanceStarting}
		running := models.AppInstanceFields{State: models.InstanceRunning}

		defaultInstanceReponses = [][]models.AppInstanceFields{
			[]models.AppInstanceFields{starting, starting},
			[]models.AppInstanceFields{starting, starting},
			[]models.AppInstanceFields{starting, running},
		}

		cmd = NewStart(ui, configRepo, mockClock, appDisplayer, appRepo, appInstancesRepo, logRepo)
	})

	runCommand := func(args ...string) {
		testcmd.RunCommand(cmd, args, requirementsFactory)
	}

	Describe("requirements", func() {
		It("fails requirements when not logged in", func() {
			runCommand("some-app-name")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})

		It("fails with usage when provided with no args", func() {
			runCommand()
			Expect(ui.FailedWithUsage).To(BeTrue())
		})
	})

	Describe("timeouts", func() {
		BeforeEach(func() {
			appRepo.UpdateAppResult = app
			appRepo.ReadReturns.App = app
			requirementsFactory.Application = app
			requirementsFactory.LoginSuccess = true
		})

		It("has sane default timeout values", func() {
			Expect(cmd.StagingTimeout).To(Equal(15 * time.Minute))
			Expect(cmd.StartupTimeout).To(Equal(5 * time.Minute))
		})

		It("can read timeout values from environment variables", func() {
			oldStaging := os.Getenv("CF_STAGING_TIMEOUT")
			oldStart := os.Getenv("CF_STARTUP_TIMEOUT")
			defer func() {
				os.Setenv("CF_STAGING_TIMEOUT", oldStaging)
				os.Setenv("CF_STARTUP_TIMEOUT", oldStart)
			}()

			os.Setenv("CF_STAGING_TIMEOUT", "6")
			os.Setenv("CF_STARTUP_TIMEOUT", "3")

			cmd = NewStart(ui, configRepo, mockClock, appDisplayer, appRepo, appInstancesRepo, logRepo)
			Expect(cmd.StagingTimeout).To(Equal(6 * time.Minute))
			Expect(cmd.StartupTimeout).To(Equal(3 * time.Minute))
		})

		Describe("when the staging timeout is zero seconds", func() {
			BeforeEach(func() {
				cmd.StagingTimeout = 0
			})

			It("can still respond to staging failures", func() {
				appInstancesRepo.GetInstancesErrorCodes = []string{"170001"}
				testcmd.RunCommand(cmd, []string{"my-app"}, requirementsFactory)

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"FAILED"},
					[]string{"Error staging app"},
				))
			})
		})
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true

			// FIXME: shouldn't need to overwrite these under test
			cmd.StagingTimeout = 50 * time.Millisecond
			cmd.StartupTimeout = 50 * time.Millisecond
			cmd.PingerThrottle = 50 * time.Millisecond
		})

		Context("when an app with the given name exists", func() {
			BeforeEach(func() {
				appRepo.UpdateAppResult = app
				appRepo.ReadReturns.App = app
				requirementsFactory.Application = app
				appInstancesRepo.GetInstancesResponses = defaultInstanceReponses
				appInstancesRepo.GetInstancesErrorCodes = defaultInstanceErrorCodes

				currentTime := time.Now()
				wrongSourceName := "DEA"
				correctSourceName := "STG"
				logRepo.TailLogMessages = []*logmessage.LogMessage{
					testlogs.NewLogMessage("Log Line 1", app.Guid, wrongSourceName, currentTime),
					testlogs.NewLogMessage("Log Line 2", app.Guid, correctSourceName, currentTime),
					testlogs.NewLogMessage("Log Line 3", app.Guid, correctSourceName, currentTime),
					testlogs.NewLogMessage("Log Line 4", app.Guid, wrongSourceName, currentTime),
				}
			})

			It("starts an app, when given the app's name", func() {
				runCommand("my-app")

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"my-app", "my-org", "my-space", "my-user"},
					[]string{"OK"},
					[]string{"0 of 2 instances running", "2 starting"},
					[]string{"started"},
				))

				Expect(requirementsFactory.ApplicationName).To(Equal("my-app"))
				Expect(appRepo.UpdateAppGuid).To(Equal("my-app-guid"))
				Expect(appDisplayer.AppToDisplay).To(Equal(app))
			})

			It("only displays staging logs when an app is starting", func() {
				runCommand("my-app")

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Log Line 2"},
					[]string{"Log Line 3"},
				))
				Expect(ui.Outputs).ToNot(ContainSubstrings(
					[]string{"Log Line 1"},
					[]string{"Log Line 4"},
				))
			})

			Context("when the app is still staging", func() {
				BeforeEach(func() {
					down := models.AppInstanceFields{State: models.InstanceDown}
					starting := models.AppInstanceFields{State: models.InstanceStarting}
					running := models.AppInstanceFields{State: models.InstanceRunning}

					appInstancesRepo.GetInstancesResponses = [][]models.AppInstanceFields{
						[]models.AppInstanceFields{},
						[]models.AppInstanceFields{},
						[]models.AppInstanceFields{down, starting},
						[]models.AppInstanceFields{starting, starting},
						[]models.AppInstanceFields{running, running},
					}
					appInstancesRepo.GetInstancesErrorCodes = []string{
						errors.APP_NOT_STAGED,
						errors.APP_NOT_STAGED,
						"", "", "",
					}

					logRepo.TailLogMessages = []*logmessage.LogMessage{
						testlogs.NewLogMessage("Log Line 1", app.Guid, LogMessageTypeStaging, time.Now()),
						testlogs.NewLogMessage("Log Line 2", app.Guid, LogMessageTypeStaging, time.Now()),
					}

				})

				It("gracefully handles starting the app", func() {
					runCommand("my-app")

					Expect(appInstancesRepo.GetInstancesAppGuid).To(Equal("my-app-guid"))
					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"Log Line 1"},
						[]string{"Log Line 2"},
						[]string{"0 of 2 instances running", "2 starting"},
					))
				})
			})

			Context("when staging the app fails", func() {
				BeforeEach(func() {
					appInstancesRepo.GetInstancesResponses = [][]models.AppInstanceFields{}
					appInstancesRepo.GetInstancesErrorCodes = []string{"170001"}
				})

				It("displays an error message when staging fails", func() {
					runCommand("my-app")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"FAILED"},
						[]string{"Error staging app"},
					))
				})
			})

			Context("when an app instance is flapping", func() {
				BeforeEach(func() {
					starting := models.AppInstanceFields{State: models.InstanceStarting}
					flapping := models.AppInstanceFields{State: models.InstanceFlapping}

					appInstancesRepo.GetInstancesResponses = [][]models.AppInstanceFields{
						[]models.AppInstanceFields{starting, starting},
						[]models.AppInstanceFields{starting, flapping},
					}
					appInstancesRepo.GetInstancesErrorCodes = []string{"", ""}

				})

				It("fails and alerts the user", func() {
					runCommand("my-app")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"0 of 2 instances running", "1 starting", "1 failing"},
						[]string{"FAILED"},
						[]string{"Start unsuccessful"},
					))
				})
			})

			Context("when waiting for the app to start times out", func() {
				var clockDestroyer chan bool

				BeforeEach(func() {
					appInstancesRepo.GetInstancesErrorCodes = []string{errors.APP_NOT_STAGED}

					clockDestroyer = make(chan bool, 1)
					go func(stopChannel chan bool) {
						for {
							select {
							case <-time.After(time.Millisecond * 100):
								mockClock.Tick()
							case <-stopChannel:
								return
							}
						}
					}(clockDestroyer)
				})

				AfterEach(func() {
					clockDestroyer <- true
				})

				It("fails and tells the user about it", func() {
					runCommand("my-app")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"FAILED"},
						[]string{"Start app timeout"},
					))
					Expect(ui.Outputs).ToNot(ContainSubstrings([]string{"instances running"}))
				})
			})

			Context("when starting the app fails", func() {
				BeforeEach(func() {
					appRepo.UpdateErr = true
				})

				It("tells the user about the failure when starting the app fails", func() {
					runCommand("my-app")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"my-app"},
						[]string{"FAILED"},
						[]string{"Error updating app."},
					))
				})
			})

			Context("when the app is already running", func() {
				BeforeEach(func() {
					requirementsFactory.Application.State = "started"
				})

				It("warns the user when the app is already running", func() {
					runCommand("my-app")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"my-app", "is already started"}),
					)
					Expect(appRepo.UpdateAppGuid).To(Equal(""))
				})
			})

			Context("when connecting to the loggregator server fails", func() {
				It("fails and tells the user", func() {
					logRepo.TailLogErr = errors.New("Ooops")
					runCommand("my-app")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"error tailing logs"},
						[]string{"Ooops"},
					))
				})
			})
		})
	})
})
