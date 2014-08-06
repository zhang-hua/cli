package broker_builder

import (
	"github.com/cloudfoundry/cli/cf/actors/service_builder"
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
	brokerRepo     api.ServiceBrokerRepository
	serviceBuilder service_builder.ServiceBuilder
}

func NewBuilder(broker api.ServiceBrokerRepository, serviceBuilder service_builder.ServiceBuilder) Builder {
	return Builder{
		brokerRepo:     broker,
		serviceBuilder: serviceBuilder,
	}
}

func (builder Builder) GetAllServiceBrokers() ([]models.ServiceBroker, error) {
	brokers := []models.ServiceBroker{}
	err := builder.brokerRepo.ListServiceBrokers(func(broker models.ServiceBroker) bool {
		brokers = append(brokers, broker)
		return true
	})

	for index, broker := range brokers {
		services, err := builder.serviceBuilder.GetServicesForBroker(broker.Guid)
		if err != nil {
			return nil, err
		}

		brokers[index].Services = services
	}
	return brokers, err
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
