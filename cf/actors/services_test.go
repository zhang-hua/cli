package actors_test

import (
	"github.com/cloudfoundry/cli/cf/actors"
	"github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	var (
		actor       actors.ServiceActor
		brokerRepo  *fakes.FakeServiceBrokerRepo
		serviceRepo *fakes.FakeServiceRepo
	)

	BeforeEach(func() {
		brokerRepo = &fakes.FakeServiceBrokerRepo{}
		serviceRepo = &fakes.FakeServiceRepo{}

		actor = actors.NewServiceActor(brokerRepo, serviceRepo)
	})

	Describe("Get Brokers with Dependencies", func() {
		Describe("Populates Dependencies", func() {
			BeforeEach(func() {
				brokerRepo.ServiceBrokers = []models.ServiceBroker{
					{Guid: "my-service-broker-guid", Name: "my-service-broker"},
					{Guid: "my-service-broker-guid2", Name: "my-service-broker2"},
				}

				serviceRepo.ListServicesFromBrokerReturns = map[string][]models.ServiceOffering{
					"my-service-broker-guid": {},
					"my-service-broker-guid2": {
						{ServiceOfferingFields: models.ServiceOfferingFields{Guid: "service-guid", Label: "my-service"}},
						{ServiceOfferingFields: models.ServiceOfferingFields{Guid: "service-guid2", Label: "my-service2"}},
					},
				}
			})

			It("Populates Services", func() {
				brokers, err := actor.GetBrokersWithDependencies()

				Expect(err).NotTo(HaveOccurred())

				Expect(len(brokers)).To(Equal(2))
				Expect(len(brokers[0].Services)).To(Equal(0))
				Expect(len(brokers[1].Services)).To(Equal(2))
				Expect(brokers[1].Services[0].Guid).To(Equal("service-guid"))
				Expect(brokers[1].Services[1].Guid).To(Equal("service-guid2"))
			})
		})
	})

})
