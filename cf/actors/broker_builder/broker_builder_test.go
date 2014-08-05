package broker_builder_test

import (
	"github.com/cloudfoundry/cli/cf/actors"
	"github.com/cloudfoundry/cli/cf/actors/broker_builder"
	"github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/models"
	. "github.com/onsi/ginkgo/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Broker Builder", func() {
	var (
		builder                   broker_builder.BrokerBuilder
		brokerRepo                *fakes.FakeServiceBrokerRepo
		serviceRepo               *fakes.FakeServiceRepo
		servicePlanRepo           *fakes.FakeServicePlanRepo
		servicePlanVisibilityRepo *fakes.FakeServicePlanVisibilityRepository
		orgRepo                   *fakes.FakeOrgRepository
	)

	BeforeEach(func() {
		brokerRepo = &fakes.FakeServiceBrokerRepo{}
		orgRepo = &fakes.FakeOrgRepository{}
		brokerBuilder = &broker_builder.FakeBrokerBuilder{}

		actor = actors.NewServiceHandler(brokerRepo, orgRepo, brokerBuilder)

		serviceBroker1 = models.ServiceBroker{Guid: "my-service-broker-guid", Name: "my-service-broker"}
		serviceBroker2 = models.ServiceBroker{Guid: "my-service-broker-guid2", Name: "my-service-broker2"}

		brokerRepo.FindByNameServiceBroker = serviceBroker2

		brokerRepo.ServiceBrokers = []models.ServiceBroker{
			serviceBroker1,
			serviceBroker2,
		}

		serviceRepo.ListServicesFromBrokerReturns = map[string][]models.ServiceOffering{
			"my-service-broker-guid": {
				{ServiceOfferingFields: models.ServiceOfferingFields{Guid: "a-guid", Label: "a-label"}},
			},
			"my-service-broker-guid2": {
				{ServiceOfferingFields: models.ServiceOfferingFields{Guid: "service-guid", Label: "my-service"}},
				{ServiceOfferingFields: models.ServiceOfferingFields{Guid: "service-guid2", Label: "my-service2"}},
			},
		}

		service2 := models.ServiceOffering{ServiceOfferingFields: models.ServiceOfferingFields{
			Label:      "my-service2",
			Guid:       "service-guid2",
			BrokerGuid: "my-service-broker-guid2"},
		}

		serviceRepo.FindServiceOfferingByLabelServiceOffering = service2

		servicePlanRepo.SearchReturns = map[string][]models.ServicePlanFields{
			"service-guid": {{Name: "service-plan", Guid: "service-plan-guid", ServiceOfferingGuid: "service-guid"},
				{Name: "other-plan", Guid: "other-plan-guid", ServiceOfferingGuid: "service-guid", Public: true}},
			"service-guid2": {{Name: "service-plan2", Guid: "service-plan2-guid", ServiceOfferingGuid: "service-guid2"}},
		}

		servicePlanVisibilityRepo.ListReturns([]models.ServicePlanVisibilityFields{
			{ServicePlanGuid: "service-plan2-guid", OrganizationGuid: "org-guid"},
			{ServicePlanGuid: "service-plan-guid", OrganizationGuid: "org-guid"},
			{ServicePlanGuid: "service-plan-guid", OrganizationGuid: "org2-guid"},
			{ServicePlanGuid: "service-plan2-guid", OrganizationGuid: "org2-guid"},
			{ServicePlanGuid: "other-plan-guid", OrganizationGuid: "org-guid"},
		}, nil)

		org1 := models.Organization{}
		org1.Name = "org1"
		org1.Guid = "org-guid"

		org2 := models.Organization{}
		org2.Name = "org2"
		org2.Guid = "org2-guid"

		orgRepo.Organizations = []models.Organization{
			org1,
			org2,
		}
	})

})
