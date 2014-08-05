package broker_builder

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/models"
)

type BrokerBuilder interface {
	GetAllServiceBrokers() ([]models.ServiceBroker, error)
	GetBrokersForServices([]models.ServiceOffering) ([]models.ServiceBroker, error)
	GetBrokerWithAllServices(brokerName string) ([]models.ServiceBroker, error)
	GetBrokerWithSpecifiedService(serviceName string) ([]models.ServiceBroker, error)
	GetSpecificBrokerForServices(string, []models.ServiceOffering) ([]models.ServiceBroker, error)
}

type Builder struct {
	brokerRepo                api.ServiceBrokerRepository
	serviceRepo               api.ServiceRepository
	servicePlanRepo           api.ServicePlanRepository
	servicePlanVisibilityRepo api.ServicePlanVisibilityRepository
	orgRepo                   api.OrganizationRepository
}

func NewBuilder(broker api.ServiceBrokerRepository, service api.ServiceRepository, plan api.ServicePlanRepository, vis api.ServicePlanVisibilityRepository, org api.OrganizationRepository) Builder {
	return Builder{
		brokerRepo:                broker,
		serviceRepo:               service,
		servicePlanRepo:           plan,
		servicePlanVisibilityRepo: vis,
		orgRepo:                   org,
	}
}

func (builder Builder) GetAllServiceBrokers() (brokers []models.ServiceBroker, err error) {
	err = builder.brokerRepo.ListServiceBrokers(func(broker models.ServiceBroker) bool {
		brokers = append(brokers, broker)
		return true
	})
	return
}

func (builder Builder) GetBrokersForServices([]models.ServiceOffering) ([]models.ServiceBroker, error) {
	return nil, nil
}

func (builder Builder) GetBrokerWithAllServices(brokerName string) ([]models.ServiceBroker, error) {
	//TEST ME!!!
	// broker, err := builder.brokerRepo.FindByName(brokerName)
	// if err != nil {
	// 	return nil, err
	// }
	//	brokers := []models.ServiceBroker{broker}
	//return builder.AttachServicesToBrokers(brokers)
	return nil, nil
}

func (builder Builder) GetBrokerWithSpecifiedService(serviceName string) ([]models.ServiceBroker, error) {
	//TEST ME!
	return nil, nil
}

func (builder Builder) GetSpecificBrokerForServices(string, []models.ServiceOffering) ([]models.ServiceBroker, error) {
	return nil, nil
}
