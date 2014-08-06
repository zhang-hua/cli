package broker_builder_test

import (
	"errors"

	"github.com/cloudfoundry/cli/cf/actors/broker_builder"
	"github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/models"

	fake_service_builder "github.com/cloudfoundry/cli/cf/actors/service_builder/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Broker Builder", func() {
	var (
		brokerBuilder broker_builder.BrokerBuilder

		serviceBuilder *fake_service_builder.FakeServiceBuilder
		brokerRepo     *fakes.FakeServiceBrokerRepo

		serviceBroker1 models.ServiceBroker

		services models.ServiceOfferings
	)

	BeforeEach(func() {
		brokerRepo = &fakes.FakeServiceBrokerRepo{}
		serviceBuilder = &fake_service_builder.FakeServiceBuilder{}
		brokerBuilder = broker_builder.NewBuilder(brokerRepo, serviceBuilder)

		serviceBroker1 = models.ServiceBroker{Guid: "my-service-broker-guid", Name: "my-service-broker"}

		publicServicePlan := models.ServicePlanFields{
			Name:   "public-service-plan",
			Guid:   "public-service-plan-guid",
			Public: true,
		}

		privateServicePlan := models.ServicePlanFields{
			Name:   "private-service-plan",
			Guid:   "private-service-plan-guid",
			Public: false,
			OrgNames: []string{
				"org-1",
				"org-2",
			},
		}

		services = models.ServiceOfferings{
			{
				ServiceOfferingFields: models.ServiceOfferingFields{
					Label: "my-public-service",
					Guid:  "my-public-service-guid",
				},
				Plans: []models.ServicePlanFields{
					publicServicePlan,
					privateServicePlan,
				},
			},
		}

		/*
			publicServicePlanVisibilityFields = models.ServicePlanVisibilityFields{
				Guid:            "public-service-plan-visibility-guid",
				ServicePlanGuid: "public-service-plan-guid",
			}

			privateServicePlanVisibilityFields = models.ServicePlanVisibilityFields{
				Guid:            "private-service-plan-visibility-guid",
				ServicePlanGuid: "private-service-plan-guid",
			}

		*/
		/*
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
		*/
	})

	Describe(".GetAllServiceBrokers", func() {
		It("returns an error if we cannot list all brokers", func() {
			brokerRepo.ListErr = true

			_, err := brokerBuilder.GetAllServiceBrokers()
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if we cannot list the services for a broker", func() {
			brokerRepo.ServiceBrokers = []models.ServiceBroker{serviceBroker1}
			serviceBuilder.GetServicesForBrokerReturns(nil, errors.New("Cannot find services"))

			_, err := brokerBuilder.GetAllServiceBrokers()
			Expect(err).To(HaveOccurred())
		})

		It("returns all service brokers populated with their services", func() {
			brokerRepo.ServiceBrokers = []models.ServiceBroker{serviceBroker1}
			serviceBuilder.GetServicesForBrokerReturns(services, nil)

			brokers, err := brokerBuilder.GetAllServiceBrokers()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(brokers)).To(Equal(1))
			Expect(brokers[0].Name).To(Equal("my-service-broker"))
			Expect(brokers[0].Services[0].Label).To(Equal("my-public-service"))
			Expect(len(brokers[0].Services[0].Plans)).To(Equal(2))
		})
	})

})
